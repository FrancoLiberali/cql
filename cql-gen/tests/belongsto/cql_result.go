// Code generated by cql-gen v0.1.0, DO NOT EDIT.
package belongsto

import preload "github.com/FrancoLiberali/cql/preload"

func (m Owned) GetOwner() (*Owner, error) {
	return preload.VerifyStructLoaded[Owner](&m.Owner)
}
