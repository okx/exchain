package rootmulti

import dbm "github.com/okex/exchain/libs/tm-db"

type metaData struct {
	version      int64
	cInfo        commitInfo
	pruneHeights []int64
	versions     []int64
}

func (m *metaData) copy() *metaData {
	ret := &metaData{}
	ret.version = m.version
	ret.cInfo = *m.cInfo.Copy()
	ret.pruneHeights = make([]int64, len(m.pruneHeights))
	copy(ret.pruneHeights, m.pruneHeights)
	ret.versions = make([]int64, len(m.versions))
	copy(ret.versions, m.versions)

	return ret
}

func (rs *Store) cacheMetadata() {
	latestVersion := getLatestVersion(rs.db)
	cInfo, err := getCommitInfo(rs.db, latestVersion)
	if err != nil && latestVersion > 0 {
		panic(latestVersion)
	}
	versions, err := getVersions(rs.db)
	if err != nil {
		panic(err)
	}
	pruneHeighs, err := getPruningHeights(rs.db, false)
	if err != nil {
		panic(err)
	}
	rs.updateMetadataToCache(&metaData{
		version:      latestVersion,
		cInfo:        cInfo,
		versions:     versions,
		pruneHeights: pruneHeighs,
	})
}

func (rs *Store) updateMetadataToCache(data *metaData) {
	rs.commitInfoVersion.Store(data.version, data.cInfo)
	rs.commitInfoVersion.Delete(data.version - MaxAsyncJob)
	rs.metadata.Store(data)
}

func (rs *Store) getMetaVersionFromCache() int64 {
	return rs.metadata.Load().(*metaData).version
}

func (rs *Store) getMetaCommitInfoFromCache() commitInfo {
	return rs.metadata.Load().(*metaData).cInfo
}

func (rs *Store) getMetaPruneHeightsFromCache() []int64 {
	return rs.metadata.Load().(*metaData).pruneHeights
}

func (rs *Store) getMetaVersions() []int64 {
	return rs.metadata.Load().(*metaData).versions
}

func (rs *Store) getCommitInfoVersion(db dbm.DB, ver int64) (commitInfo, error) {
	cInfo, ok := rs.commitInfoVersion.Load(ver)
	if !ok {
		return getCommitInfo(db, ver)
	}

	return cInfo.(commitInfo), nil
}

func (rs *Store) getLatestVersion() int64 {
	if rs.enableAsyncJob {
		return rs.getMetaVersionFromCache()
	}
	return getLatestVersion(rs.db)
}

func (rs *Store) getVersions() ([]int64, error) {
	if rs.enableAsyncJob {
		return rs.getMetaVersions(), nil
	}
	return getVersions(rs.db)
}

func (rs *Store) getPruningHeights() ([]int64, error) {
	if rs.enableAsyncJob {
		return rs.getMetaPruneHeightsFromCache(), nil
	}
	return getPruningHeights(rs.db, false)
}
