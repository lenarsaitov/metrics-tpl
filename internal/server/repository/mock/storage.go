package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	models "github.com/lenarsaitov/metrics-tpl/internal/server/models"
	reflect "reflect"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// AddCounter mocks base method.
func (m *MockStorage) AddCounter(ctx context.Context, name string, value int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCounter", ctx, name, value)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddCounter indicates an expected call of AddCounter.
func (mr *MockStorageMockRecorder) AddCounter(ctx, name, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCounter", reflect.TypeOf((*MockStorage)(nil).AddCounter), ctx, name, value)
}

// GetAll mocks base method.
func (m *MockStorage) GetAll(ctx context.Context) (models.Metrics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].(models.Metrics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockStorageMockRecorder) GetAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockStorage)(nil).GetAll), ctx)
}

// GetCounterMetric mocks base method.
func (m *MockStorage) GetCounterMetric(ctx context.Context, name string) (*int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCounterMetric", ctx, name)
	ret0, _ := ret[0].(*int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCounterMetric indicates an expected call of GetCounterMetric.
func (mr *MockStorageMockRecorder) GetCounterMetric(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCounterMetric", reflect.TypeOf((*MockStorage)(nil).GetCounterMetric), ctx, name)
}

// GetGaugeMetric mocks base method.
func (m *MockStorage) GetGaugeMetric(ctx context.Context, name string) (*float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGaugeMetric", ctx, name)
	ret0, _ := ret[0].(*float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGaugeMetric indicates an expected call of GetGaugeMetric.
func (mr *MockStorageMockRecorder) GetGaugeMetric(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGaugeMetric", reflect.TypeOf((*MockStorage)(nil).GetGaugeMetric), ctx, name)
}

// ReplaceGauge mocks base method.
func (m *MockStorage) ReplaceGauge(ctx context.Context, name string, value float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplaceGauge", ctx, name, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReplaceGauge indicates an expected call of ReplaceGauge.
func (mr *MockStorageMockRecorder) ReplaceGauge(ctx, name, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplaceGauge", reflect.TypeOf((*MockStorage)(nil).ReplaceGauge), ctx, name, value)
}
