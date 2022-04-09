package app

import (
	"sort"

	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	upgradetypes "github.com/okex/exchain/libs/cosmos-sdk/types/upgrade"
	"github.com/okex/exchain/libs/cosmos-sdk/x/params"
	"github.com/okex/exchain/libs/cosmos-sdk/x/params/subspace"
)

func (app *OKExChainApp) setupUpgradeModules() {
	heightTasks, paramMap, pip, prunePip := app.CollectUpgradeModules(app.mm)

	app.heightTasks = heightTasks

	if pip != nil {
		app.GetCMS().SetPruneHeightFilterPipeline(prunePip)
		app.GetCMS().SetCommitHeightFilterPipeline(pip)
	}

	vs := app.subspaces
	for k, vv := range paramMap {
		supace, exist := vs[k]
		if !exist {
			continue
		}
		vs[k] = supace.LazyWithKeyTable(subspace.NewKeyTable(vv.ParamSetPairs()...))
	}
}

func (o *OKExChainApp) CollectUpgradeModules(m *module.Manager) (map[int64]*upgradetypes.HeightTasks, map[string]params.ParamSet, types.HeightFilterPipeline, types.HeightFilterPipeline) {
	hm := make(map[int64]*upgradetypes.HeightTasks)
	hStoreInfoModule := make(map[int64]map[string]struct{})
	paramsRet := make(map[string]params.ParamSet)
	for _, mm := range m.Modules {
		if ada, ok := mm.(upgradetypes.UpgradeModule); ok {
			set := ada.RegisterParam()
			if set != nil {
				if _, exist := paramsRet[ada.ModuleName()]; !exist {
					paramsRet[ada.ModuleName()] = set
				}
			}
			h := ada.UpgradeHeight()
			if h > 0 {
				h++
			}
			storeInfoModule := hStoreInfoModule[h]
			if storeInfoModule == nil {
				storeInfoModule = make(map[string]struct{})
				hStoreInfoModule[h] = storeInfoModule
			}
			names := ada.BlockStoreModules()
			for _, n := range names {
				storeInfoModule[n] = struct{}{}
			}
			t := ada.RegisterTask()
			if t == nil {
				continue
			}
			if err := t.ValidateBasic(); nil != err {
				panic(err)
			}
			taskList := hm[h]
			if taskList == nil {
				v := make(upgradetypes.HeightTasks, 0)
				taskList = &v
				hm[h] = taskList
			}
			*taskList = append(*taskList, t)
		}
	}

	for _, v := range hm {
		sort.Sort(*v)
	}

	commitPip, prunePip := collectStorePipeline(hStoreInfoModule)

	return hm, paramsRet, commitPip, prunePip
}

func collectStorePipeline(hStoreInfoModule map[int64]map[string]struct{}) (types.HeightFilterPipeline, types.HeightFilterPipeline) {
	var (
		pip      types.HeightFilterPipeline
		prunePip types.HeightFilterPipeline
	)

	for storeH, storeMap := range hStoreInfoModule {
		filterM := copyBlockStoreMap(storeMap)
		if storeH < 0 {
			continue
		}
		hh := storeH
		height := hh - 1
		// filter block module
		blockModuleFilter := func(str string) bool {
			_, exist := filterM[str]
			return exist
		}

		commitF := func(h int64) func(str string) bool {
			if hh == 0 {
				return blockModuleFilter
			}
			if h >= height {
				// call next filter
				return nil
			}
			return blockModuleFilter
		}
		pruneF := func(h int64) func(str string) bool {
			if hh == 0 {
				return blockModuleFilter
			}
			// note: prune's version  > commit version,thus the condition will be '>' rather than '>='
			if h > height {
				// call next filter
				return nil
			}
			return blockModuleFilter
		}

		pip = linkPipeline(pip, commitF)
		prunePip = linkPipeline(prunePip, pruneF)
	}

	return pip, prunePip
}

func copyBlockStoreMap(m map[string]struct{}) map[string]struct{} {
	ret := make(map[string]struct{})
	for k, _ := range m {
		ret[k] = struct{}{}
	}
	return ret
}

func linkPipeline(p types.HeightFilterPipeline, f func(h int64) func(str string) bool) types.HeightFilterPipeline {
	if p == nil {
		p = f
	} else {
		p = types.LinkPipeline(f, p)
	}
	return p
}
