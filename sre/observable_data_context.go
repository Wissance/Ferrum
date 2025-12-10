package sre

import (
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/managers"
	sf "github.com/wissance/stringFormatter"
)

const realmKeyLabelTemplate = "{0}_{1}_{2}"
const realmRelateResourceKeyWithoutId = "{0}_{1}_for_realm_{2}"
const realmRelatedKeyLabelTemplate = "{0}_{1}_{2}_for_realm_{3}"
const realmLessKeyLabelTemplate = "{0}_{1}"

// ObservableDataContext struct that implements managers.DataContext
type ObservableDataContext struct {
	metricsCollector *MetricsCollector
	dataContext      *managers.DataContext
}

// CreateObservableDataContext function that wraps DataContext operations for observability
func CreateObservableDataContext(metricsCollector *MetricsCollector, dataContext *managers.DataContext) managers.DataContext {
	return &ObservableDataContext{
		metricsCollector: metricsCollector,
		dataContext:      dataContext,
	}
}

// IsAvailable this function we do not observe yet
func (mn *ObservableDataContext) IsAvailable() bool {
	return mn.IsAvailable()
}

// GetRealm function for Observe (SRE) getting realm operation
/* Function that wraps managers.DataContext GetRealm with getting the following metrics:
 *   1. DataSource requests duration - DataSourceRequestDurations in us
 *   2. DataSourceRequestsTotalCount - count of requests with key and status labels
 */
func (mn *ObservableDataContext) GetRealm(realmName string) (*data.Realm, error) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(us)
	}))
	r, e := mn.GetRealm(realmName)
	timer.ObserveDuration()
	key := sf.Format(realmKeyLabelTemplate, "get", data.REALM, realmName)
	if e == nil {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, successStatus).Inc()
	} else {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, failureStatus).Inc()
	}
	return r, e
}

// GetUsers function for getting all Realm User
/* Function that wraps managers.DataContext GetUsers with getting the following metrics:
 *   1. DataSource requests duration - DataSourceRequestDurations in us
 *   2. DataSourceRequestsTotalCount - count of requests with key and status labels
 */
func (mn *ObservableDataContext) GetUsers(realmName string) ([]data.User, error) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(us)
	}))
	u, e := mn.GetUsers(realmName)
	timer.ObserveDuration()
	key := sf.Format(realmRelateResourceKeyWithoutId, "getmany", data.USER, realmName)
	if e == nil {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, successStatus).Inc()
	} else {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, failureStatus).Inc()
	}
	return u, e
}

// GetClient function for getting Realm Client by name
/* Function that wraps managers.DataContext GetClient with getting the following metrics:
 *   1. DataSource requests duration - DataSourceRequestDurations in us
 *   2. DataSourceRequestsTotalCount - count of requests with key and status labels
 */
func (mn *ObservableDataContext) GetClient(realmName string, clientName string) (*data.Client, error) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(us)
	}))
	c, e := mn.GetClient(realmName, clientName)
	timer.ObserveDuration()
	key := sf.Format(realmRelatedKeyLabelTemplate, "get", data.CLIENT, clientName, realmName)
	if e == nil {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, successStatus).Inc()
	} else {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, failureStatus).Inc()
	}
	return c, e
}

// GetUser function for getting Realm User by userName
/* Function that wraps managers.DataContext GetUser with getting the following metrics:
 *   1. DataSource requests duration - DataSourceRequestDurations in us
 *   2. DataSourceRequestsTotalCount - count of requests with key and status labels
 */
func (mn *ObservableDataContext) GetUser(realmName string, userName string) (data.User, error) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(us)
	}))
	u, e := mn.GetUser(realmName, userName)
	timer.ObserveDuration()
	key := sf.Format(realmRelatedKeyLabelTemplate, "get", data.USER, userName, realmName)
	if e == nil {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, successStatus).Inc()
	} else {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, failureStatus).Inc()
	}
	return u, e
}

// GetUserById function for getting Realm User by UserId (uuid)
/* Function that wraps managers.DataContext GetUserById with getting the following metrics:
 *   1. DataSource requests duration - DataSourceRequestDurations in us
 *   2. DataSourceRequestsTotalCount - count of requests with key and status labels
 */
func (mn *ObservableDataContext) GetUserById(realmName string, userId uuid.UUID) (data.User, error) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(us)
	}))
	u, e := mn.GetUserById(realmName, userId)
	timer.ObserveDuration()
	key := sf.Format(realmRelatedKeyLabelTemplate, "get", data.USER, userId.String(), realmName)
	if e == nil {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, successStatus).Inc()
	} else {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, failureStatus).Inc()
	}
	return u, e
}

// CreateRealm creates new data.Realm in a data store, receive realmData unmarshalled json in a data.Realm
/* Function that wraps managers.DataContext CreateRealm with getting the following metrics:
 *   1. DataSource requests duration - DataSourceRequestDurations in us
 *   2. DataSourceRequestsTotalCount - count of requests with key and status labels
 */
func (mn *ObservableDataContext) CreateRealm(realmData data.Realm) error {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(us)
	}))
	e := mn.CreateRealm(realmData)
	timer.ObserveDuration()
	key := sf.Format(realmKeyLabelTemplate, "create", data.REALM, realmData.Name)
	if e == nil {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, successStatus).Inc()
	} else {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, failureStatus).Inc()
	}

	return e
}

// CreateClient creates new data.Client in a data store, requires to pass realmName (because client name is not unique), clientData is an unmarshalled json of type data.Client
/* Function that wraps managers.DataContext CreateClient with getting the following metrics:
 *   1. DataSource requests duration - DataSourceRequestDurations in us
 *   2. DataSourceRequestsTotalCount - count of requests with key and status labels
 */
func (mn *ObservableDataContext) CreateClient(realmName string, clientData data.Client) error {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(us)
	}))
	timer.ObserveDuration()
	e := mn.CreateClient(realmName, clientData)
	key := sf.Format(realmRelatedKeyLabelTemplate, "create", data.CLIENT, clientData.Name, realmName)
	if e == nil {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, successStatus).Inc()
	} else {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, failureStatus).Inc()
	}

	return e
}

// CreateUser creates new data.User in a data store within a realm with name = realmName
/* Function that wraps managers.DataContext CreateUser with getting the following metrics:
 *   1. DataSource requests duration - DataSourceRequestDurations in us
 *   2. DataSourceRequestsTotalCount - count of requests with key and status labels
 */
func (mn *ObservableDataContext) CreateUser(realmName string, userData data.User) error {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(us)
	}))
	e := mn.CreateUser(realmName, userData)
	timer.ObserveDuration()
	key := sf.Format(realmRelatedKeyLabelTemplate, "create", data.USER, userData.GetUsername(), realmName)
	if e == nil {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, successStatus).Inc()
	} else {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, failureStatus).Inc()
	}
	return e
}

// UpdateRealm updates existing data.Realm in a data store within name = realmData, and new data = realmData
/* Function that wraps managers.DataContext UpdateRealm with getting the following metrics:
 *   1. DataSource requests duration - DataSourceRequestDurations in us
 *   2. DataSourceRequestsTotalCount - count of requests with key and status labels
 */
func (mn *ObservableDataContext) UpdateRealm(realmName string, realmData data.Realm) error {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(us)
	}))
	e := mn.UpdateRealm(realmName, realmData)
	timer.ObserveDuration()
	key := sf.Format(realmKeyLabelTemplate, "update", data.REALM, realmName)
	if e == nil {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, successStatus).Inc()
	} else {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, failureStatus).Inc()
	}
	return e
}

// UpdateClient updates existing data.Client in a data store with name = clientName and new data = clientData
/* Function that wraps managers.DataContext UpdateClient with getting the following metrics:
 *   1. DataSource requests duration - DataSourceRequestDurations in us
 *   2. DataSourceRequestsTotalCount - count of requests with key and status labels
 */
func (mn *ObservableDataContext) UpdateClient(realmName string, clientName string, clientData data.Client) error {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(us)
	}))
	e := mn.UpdateClient(realmName, clientName, clientData)
	timer.ObserveDuration()
	key := sf.Format(realmRelatedKeyLabelTemplate, "update", data.CLIENT, clientName, realmName)
	if e == nil {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, successStatus).Inc()
	} else {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, failureStatus).Inc()
	}
	return e
}

// UpdateUser updates existing data.User in a data store with realm name = realName, username = userName and data=userData
/* Function that wraps managers.DataContext UpdateUser with getting the following metrics:
 *   1. DataSource requests duration - DataSourceRequestDurations in us
 *   2. DataSourceRequestsTotalCount - count of requests with key and status labels
 */
func (mn *ObservableDataContext) UpdateUser(realmName string, userName string, userData data.User) error {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(us)
	}))
	e := mn.UpdateUser(realmName, userName, userData)
	timer.ObserveDuration()
	key := sf.Format(realmRelatedKeyLabelTemplate, "update", data.USER, userName, realmName)
	if e == nil {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, successStatus).Inc()
	} else {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, failureStatus).Inc()
	}
	return e
}

// DeleteRealm removes realm from data storage (Should be a CASCADE remove of all related Users and Clients)
/* Function that wraps managers.DataContext DeleteRealm with getting the following metrics:
 *   1. DataSource requests duration - DataSourceRequestDurations in us
 *   2. DataSourceRequestsTotalCount - count of requests with key and status labels
 */
func (mn *ObservableDataContext) DeleteRealm(realmName string) error {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(us)
	}))
	e := mn.DeleteRealm(realmName)
	timer.ObserveDuration()
	key := sf.Format(realmKeyLabelTemplate, "delete", data.REALM, realmName)
	if e == nil {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, successStatus).Inc()
	} else {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, failureStatus).Inc()
	}
	return e
}

// DeleteClient removes client with name = clientName from realm with name = clientName
/* Function that wraps managers.DataContext DeleteClient with getting the following metrics:
 *   1. DataSource requests duration - DataSourceRequestDurations in us
 *   2. DataSourceRequestsTotalCount - count of requests with key and status labels
 */
func (mn *ObservableDataContext) DeleteClient(realmName string, clientName string) error {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(us)
	}))
	e := mn.DeleteClient(realmName, clientName)
	timer.ObserveDuration()
	key := sf.Format(realmRelatedKeyLabelTemplate, "delete", data.CLIENT, clientName, realmName)
	if e == nil {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, successStatus).Inc()
	} else {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, failureStatus).Inc()
	}
	return e
}

// DeleteUser removes data.User from data store by user (userName) and realm (realmName) name respectively
/* Function that wraps managers.DataContext DeleteUser with getting the following metrics:
 *   1. DataSource requests duration - DataSourceRequestDurations in us
 *   2. DataSourceRequestsTotalCount - count of requests with key and status labels
 */
func (mn *ObservableDataContext) DeleteUser(realmName string, userName string) error {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(us)
	}))
	e := mn.DeleteUser(realmName, userName)
	timer.ObserveDuration()
	key := sf.Format(realmRelatedKeyLabelTemplate, "delete", data.USER, userName, realmName)
	if e == nil {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, successStatus).Inc()
	} else {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, failureStatus).Inc()
	}
	return e
}

// GetUserFederationConfig function for getting data.UserFederationServiceConfig from data store by realm (realmName) name respectively
/* Function that wraps managers.DataContext GetUserFederationConfig with getting the following metrics:
 *   1. DataSource requests duration - DataSourceRequestDurations in us
 *   2. DataSourceRequestsTotalCount - count of requests with key and status labels
 */
func (mn *ObservableDataContext) GetUserFederationConfig(realmName string, configName string) (*data.UserFederationServiceConfig, error) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(us)
	}))
	c, e := mn.GetUserFederationConfig(realmName, configName)
	timer.ObserveDuration()
	key := sf.Format(realmRelatedKeyLabelTemplate, "get", data.USER_FEDERATION_SERVICE_CONFIG, configName, realmName)
	if e == nil {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, successStatus).Inc()
	} else {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, failureStatus).Inc()
	}
	return c, e
}

// CreateUserFederationConfig function for creating data.UserFederationServiceConfig in data.Realm
/* Function that wraps managers.DataContext CreateUserFederationConfig with getting the following metrics:
 *   1. DataSource requests duration - DataSourceRequestDurations in us
 *   2. DataSourceRequestsTotalCount - count of requests with key and status labels
 */
func (mn *ObservableDataContext) CreateUserFederationConfig(realmName string, userFederationConfig data.UserFederationServiceConfig) error {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(us)
	}))
	e := mn.CreateUserFederationConfig(realmName, userFederationConfig)
	timer.ObserveDuration()
	key := sf.Format(realmRelatedKeyLabelTemplate, "create", data.USER_FEDERATION_SERVICE_CONFIG, userFederationConfig.Name, realmName)
	if e == nil {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, successStatus).Inc()
	} else {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, failureStatus).Inc()
	}
	return e
}

// UpdateUserFederationConfig function for updating the data.UserFederationServiceConfig from data store by config (configName) and realm (realmName) name respectively
/* Function that wraps managers.DataContext UpdateUserFederationConfig with getting the following metrics:
 *   1. DataSource requests duration - DataSourceRequestDurations in us
 *   2. DataSourceRequestsTotalCount - count of requests with key and status labels
 */
func (mn *ObservableDataContext) UpdateUserFederationConfig(realmName string, configName string, userFederationConfig data.UserFederationServiceConfig) error {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(us)
	}))
	e := mn.UpdateUserFederationConfig(realmName, configName, userFederationConfig)
	timer.ObserveDuration()
	key := sf.Format(realmRelatedKeyLabelTemplate, "update", data.USER_FEDERATION_SERVICE_CONFIG, configName, realmName)
	if e == nil {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, successStatus).Inc()
	} else {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, failureStatus).Inc()
	}
	return e
}

// DeleteUserFederationConfig function for removing data.UserFederationServiceConfig from data store by realm (realmName) name respectively
/* Function that wraps managers.DataContext DeleteUserFederationConfig with getting the following metrics:
 *   1. DataSource requests duration - DataSourceRequestDurations in us
 *   2. DataSourceRequestsTotalCount - count of requests with key and status labels
 */
func (mn *ObservableDataContext) DeleteUserFederationConfig(realmName string, configName string) error {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(us)
	}))
	e := mn.DeleteUserFederationConfig(realmName, configName)
	timer.ObserveDuration()
	key := sf.Format(realmRelatedKeyLabelTemplate, "delete", data.USER_FEDERATION_SERVICE_CONFIG, configName, realmName)
	if e == nil {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, successStatus).Inc()
	} else {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, failureStatus).Inc()
	}
	return e
}

// GetServerSettings function that returns ServerSettings
/* Function that wraps managers.DataContext GetServerSettings with getting the following metrics:
 *   1. DataSource requests duration - DataSourceRequestDurations in us
 *   2. DataSourceRequestsTotalCount - count of requests with key and status labels
 */
func (mn *ObservableDataContext) GetServerSettings() (*data.ServerSettings, error) {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(us)
	}))
	s, e := mn.GetServerSettings()
	timer.ObserveDuration()
	key := sf.Format(realmLessKeyLabelTemplate, "get", data.USER_FEDERATION_SERVICE_CONFIG)
	if e == nil {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, successStatus).Inc()
	} else {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, failureStatus).Inc()
	}
	return s, e
}

// SetServerSettings function that updates ServerSettings by full new settings replace
/* Function that wraps managers.DataContext SetServerSettings with getting the following metrics:
 *   1. DataSource requests duration - DataSourceRequestDurations in us
 *   2. DataSourceRequestsTotalCount - count of requests with key and status labels
 */
func (mn *ObservableDataContext) SetServerSettings(settings *data.ServerSettings) error {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(us)
	}))
	e := mn.SetServerSettings(settings)
	timer.ObserveDuration()
	key := sf.Format(realmLessKeyLabelTemplate, "set", data.USER_FEDERATION_SERVICE_CONFIG)
	if e == nil {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, successStatus).Inc()
	} else {
		mn.metricsCollector.DataSourceRequestsTotalCount.WithLabelValues(key, failureStatus).Inc()
	}
	return e
}
