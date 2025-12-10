package sre

import (
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/wissance/Ferrum/data"
	"github.com/wissance/Ferrum/managers"
	sf "github.com/wissance/stringFormatter"
)

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
		us := v * 1000000 // make milliseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(us)
	}))
	r, e := mn.GetRealm(realmName)
	timer.ObserveDuration()
	key := sf.Format("{0}_{1}_{2}", "get", data.REALM, realmName)
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
		ms := v * 1000 // make milliseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(ms)
	}))
	u, e := mn.GetUsers(realmName)
	timer.ObserveDuration()
	key := sf.Format("{0}_{1}_for_realm_{2}", "getmany", data.USER, realmName)
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
		ms := v * 1000 // make milliseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(ms)
	}))
	c, e := mn.GetClient(realmName, clientName)
	timer.ObserveDuration()
	key := sf.Format("{0}_{1}_{2}_for_realm_{3}", "get", data.CLIENT, clientName, realmName)
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
		ms := v * 1000 // make milliseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(ms)
	}))
	u, e := mn.GetUser(realmName, userName)
	timer.ObserveDuration()
	key := sf.Format("{0}_{1}_{2}_for_realm_{3}", "get", data.USER, userName, realmName)
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
		ms := v * 1000 // make milliseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(ms)
	}))
	u, e := mn.GetUserById(realmName, userId)
	timer.ObserveDuration()
	key := sf.Format("{0}_{1}_{2}_for_realm_{3}", "get", data.USER, userId.String(), realmName)
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
		ms := v * 1000 // make milliseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(ms)
	}))
	e := mn.CreateRealm(realmData)
	timer.ObserveDuration()
	key := sf.Format("{0}_{1}_{2}", "create", data.REALM, realmData.Name)
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
		ms := v * 1000 // make milliseconds
		mn.metricsCollector.DataSourceRequestDurations.Observe(ms)
	}))
	timer.ObserveDuration()
	e := mn.CreateClient(realmName, clientData)
	key := sf.Format("{0}_{1}_{2}_for_realm_{3}", "create", data.CLIENT, clientData.Name, realmName)
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
	return mn.CreateUser(realmName, userData)
}

// UpdateRealm updates existing data.Realm in a data store within name = realmData, and new data = realmData
func (mn *ObservableDataContext) UpdateRealm(realmName string, realmData data.Realm) error {
	return mn.UpdateRealm(realmName, realmData)
}

// UpdateClient updates existing data.Client in a data store with name = clientName and new data = clientData
func (mn *ObservableDataContext) UpdateClient(realmName string, clientName string, clientData data.Client) error {
	return mn.UpdateClient(realmName, clientName, clientData)
}

// UpdateUser updates existing data.User in a data store with realm name = realName, username = userName and data=userData
func (mn *ObservableDataContext) UpdateUser(realmName string, userName string, userData data.User) error {
	return mn.UpdateUser(realmName, userName, userData)
}

// DeleteRealm removes realm from data storage (Should be a CASCADE remove of all related Users and Clients)
func (mn *ObservableDataContext) DeleteRealm(realmName string) error {
	return mn.DeleteRealm(realmName)
}

// DeleteClient removes client with name = clientName from realm with name = clientName
func (mn *ObservableDataContext) DeleteClient(realmName string, clientName string) error {
	return mn.DeleteClient(realmName, clientName)
}

// DeleteUser removes data.User from data store by user (userName) and realm (realmName) name respectively
func (mn *ObservableDataContext) DeleteUser(realmName string, userName string) error {
	return mn.DeleteUser(realmName, userName)
}

func (mn *ObservableDataContext) GetUserFederationConfig(realmName string, configName string) (*data.UserFederationServiceConfig, error) {
	return mn.GetUserFederationConfig(realmName, configName)
}

func (mn *ObservableDataContext) CreateUserFederationConfig(realmName string, userFederationConfig data.UserFederationServiceConfig) error {
	return mn.CreateUserFederationConfig(realmName, userFederationConfig)
}

func (mn *ObservableDataContext) UpdateUserFederationConfig(realmName string, configName string, userFederationConfig data.UserFederationServiceConfig) error {
	return mn.UpdateUserFederationConfig(realmName, configName, userFederationConfig)
}

func (mn *ObservableDataContext) DeleteUserFederationConfig(realmName string, configName string) error {
	return mn.DeleteUserFederationConfig(realmName, configName)
}

// GetServerSettings function that returns ServerSettings
func (mn *ObservableDataContext) GetServerSettings() (*data.ServerSettings, error) {
	return mn.GetServerSettings()
}

// SetServerSettings function that updates ServerSettings by full new settings replace
func (mn *ObservableDataContext) SetServerSettings(settings *data.ServerSettings) error {
	return mn.SetServerSettings(settings)
}
