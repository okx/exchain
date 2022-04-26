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
	heightTasks, paramMap, cf, pf, vf := app.CollectUpgradeModules(app.mm)

	app.heightTasks = heightTasks

	app.GetCMS().AppendCommitFilters(cf)
	app.GetCMS().AppendPruneFilters(pf)
	app.GetCMS().AppendVersionFilters(vf)

	vs := app.subspaces
	for k, vv := range paramMap {
		supace, exist := vs[k]
		if !exist {
			continue
		}
		vs[k] = supace.LazyWithKeyTable(subspace.NewKeyTable(vv.ParamSetPairs()...))
	}
}

func (o *OKExChainApp) CollectUpgradeModules(m *module.Manager) (map[int64]*upgradetypes.HeightTasks,
	map[string]params.ParamSet, []types.StoreFilter, []types.StoreFilter, []types.VersionFilter) {
	hm := make(map[int64]*upgradetypes.HeightTasks)
	//hStoreInfoModule := make(map[int64]map[string]upgradetypes.HandleStore)
	paramsRet := make(map[string]params.ParamSet)
	commitFiltreMap := make(map[*types.StoreFilter]struct{})
	pruneFilterMap := make(map[*types.StoreFilter]struct{})
	versionFilterMap := make(map[*types.VersionFilter]struct{})

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

			cf := ada.CommitFilter()
			if cf != nil {
				if _, exist := commitFiltreMap[cf]; !exist {
					commitFiltreMap[cf] = struct{}{}
				}
			}
			pf := ada.PruneFilter()
			if pf != nil {
				if _, exist := pruneFilterMap[pf]; !exist {
					pruneFilterMap[pf] = struct{}{}
				}
			}
			vf := ada.VersionFilter()
			if vf != nil {
				if _, exist := versionFilterMap[vf]; !exist {
					versionFilterMap[vf] = struct{}{}
				}
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

	commitFilters := make([]types.StoreFilter, 0)
	pruneFilters := make([]types.StoreFilter, 0)
	versionFilters := make([]types.VersionFilter, 0)
	for pointerFilter, _ := range commitFiltreMap {
		commitFilters = append(commitFilters, *pointerFilter)
	}
	for pointerFilter, _ := range pruneFilterMap {
		pruneFilters = append(pruneFilters, *pointerFilter)
	}
	for pointerFilter, _ := range versionFilterMap {
		versionFilters = append(versionFilters, *pointerFilter)
	}

	//commitPip, prunePip, versionPip := collectStorePipeline(hStoreInfoModule)

	return hm, paramsRet, commitFilters, pruneFilters, versionFilters
}

func collectStorePipeline(hStoreInfoModule map[int64]map[string]upgradetypes.HandleStore) (types.HeightFilterPipeline, types.PrunePipeline, types.VersionFilterPipeline) {
	var (
		pip        types.HeightFilterPipeline
		prunePip   types.PrunePipeline
		versionPip types.VersionFilterPipeline
	)

	for storeH, storeMap := range hStoreInfoModule {
		if storeH < 0 {
			continue
		}
		filterM := copyBlockStoreMap(storeMap)
		hh := storeH // storeH: upgardeHeight+1
		height := hh - 1
		// filter block module
		blockModuleFilter := func(str string) bool {
			_, exist := filterM[str]
			return exist
		}

		commitF := func(h int64) func(_ string, st types.CommitKVStore) bool {
			if hh == 0 || h < height {
				return func(str string, _ types.CommitKVStore) bool {
					return blockModuleFilter(str)
				}
			}
			if h == height {
				return func(str string, st types.CommitKVStore) bool {
					handler := filterM[str]
					if nil != handler && nil != st {
						handler(st, height)
					}
					return false
				}
			}
			// call next filter
			return func(_ string, _ types.CommitKVStore) bool {
				return false
			}
		}
		pruneF := func(h int64) func(str string) bool {
			if hh == 0 {
				return blockModuleFilter
			}
			// note: prune's version  > commit version,thus the condition will be '>' rather than '>='
			if h > height {
				// call next filter
				return func(_ string) bool {
					return false
				}
			}
			return blockModuleFilter
		}
		versionF := func(h int64) func(cb func(string, int64)) {
			//if h < height {
			//	return nil
			//}
			if h < 0 {
				return func(cb func(string, int64)) {}
			}

			return func(cb func(name string, version int64)) {

				for k, _ := range filterM {
					cb(k, hh-1)
				}
			}
		}

		pip = linkPipeline(pip, commitF)
		prunePip = linkPrunePipeline(prunePip, pruneF)
		versionPip = linkPipeline2(versionPip, versionF)
	}

	return pip, prunePip, versionPip
}

func copyBlockStoreMap(m map[string]upgradetypes.HandleStore) map[string]upgradetypes.HandleStore {
	ret := make(map[string]upgradetypes.HandleStore)
	for k, v := range m {
		ret[k] = v
	}
	return ret
}

func linkPipeline(p types.HeightFilterPipeline, f types.HeightFilterPipeline) types.HeightFilterPipeline {
	if p == nil {
		p = f
	} else {
		p = types.LinkPipeline(f, p)
	}
	return p
}
func linkPrunePipeline(p types.PrunePipeline, f types.PrunePipeline) types.PrunePipeline {
	if p == nil {
		p = f
	} else {
		p = types.LinkPrunePipeline(f, p)
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
