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
	heightTasks, paramMap, pip, prunePip, versionPip := app.CollectUpgradeModules(app.mm)

	app.heightTasks = heightTasks

	if pip != nil {
		app.GetCMS().SetPruneHeightFilterPipeline(prunePip)
		app.GetCMS().SetCommitHeightFilterPipeline(pip)
		app.GetCMS().SetVersionFilterPipeline(versionPip)
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

func (o *OKExChainApp) CollectUpgradeModules(m *module.Manager) (map[int64]*upgradetypes.HeightTasks, map[string]params.ParamSet, types.HeightFilterPipeline, types.HeightFilterPipeline, types.VersionFilterPipeline) {
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

	commitPip, prunePip, versionPip := collectStorePipeline(hStoreInfoModule)

	return hm, paramsRet, commitPip, prunePip, versionPip
}

func collectStorePipeline(hStoreInfoModule map[int64]map[string]struct{}) (types.HeightFilterPipeline, types.HeightFilterPipeline, types.VersionFilterPipeline) {
	var (
		pip        types.HeightFilterPipeline
		prunePip   types.HeightFilterPipeline
		versionPip types.VersionFilterPipeline
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
		versionF := func(h int64) func(cb func(string, int64)) {
			//if h < height {
			//	return nil
			//}
			if h < 0 {
				return nil
			}

			return func(cb func(name string, version int64)) {

				for k, _ := range filterM {
					cb(k, hh-1)
				}
			}
		}

		pip = linkPipeline(pip, commitF)
		prunePip = linkPipeline(prunePip, pruneF)
		versionPip = linkPipeline2(versionPip, versionF)
	}

	return pip, prunePip, versionPip
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

func linkPipeline2(p types.VersionFilterPipeline, f func(h int64) func(func(string, int64))) types.VersionFilterPipeline {
	if p == nil {
		p = f
	} else {
		p = types.LinkPipeline2(f, p)
	}
	return p
}
