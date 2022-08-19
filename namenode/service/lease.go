package NNService

import (
	"time"
)

func NewLeaseMgr() *LeaseMgr {
	leaseMgr := LeaseMgr{
		FileToMetaMap: make(map[string]Lease),
	}

	go leaseMgr.Monitor()
	return &leaseMgr
}

func (leaseMgr *LeaseMgr) Monitor() {
	delay := time.Minute * 3
	time.Sleep(delay)
	leaseMgr.mu.Lock()
	defer leaseMgr.mu.Unlock()

	for f, meta := range leaseMgr.FileToMetaMap {
		// TODO: HardLimit
		if time.Since(time.Unix(meta.LastUpdate, 0)) > delay {
			leaseMgr.Revoke(meta.Holder, f)
		}
	}
	leaseMgr.Monitor()
}

func (leaseMgr *LeaseMgr) Grant(client string, f string) bool {
	leaseMgr.mu.Lock()
	defer leaseMgr.mu.Unlock()

	_, present := leaseMgr.FileToMetaMap[f]
	if present {
		return false
	}

	meta := Lease{
		Holder:     client,
		LastUpdate: time.Now().Unix(),
	}

	leaseMgr.FileToMetaMap[f] = meta
	return true
}

func (leaseMgr *LeaseMgr) HasLock(f string) bool {
	leaseMgr.mu.Lock()
	defer leaseMgr.mu.Unlock()

	// TODO: SoftLimit
	_, present := leaseMgr.FileToMetaMap[f]
	return present
}

func (leaseMgr *LeaseMgr) Revoke(client string, f string) {
	leaseMgr.mu.Lock()
	defer leaseMgr.mu.Unlock()
	delete(leaseMgr.FileToMetaMap, f)
}

func (leaseMgr *LeaseMgr) Renew(client string, f string) bool {
	leaseMgr.mu.Lock()
	defer leaseMgr.mu.Unlock()

	meta, present := leaseMgr.FileToMetaMap[f]
	if present {
		if meta.Holder == client {
			meta := Lease{
				Holder:     client,
				LastUpdate: time.Now().Unix(),
			}
			leaseMgr.FileToMetaMap[f] = meta
			return true
		}
	}
	return false
}
