package common

import "store/common/datastore"

type Env struct {
	DB          datastore.DataStore
	AccessToken string
	Debug       bool
}
