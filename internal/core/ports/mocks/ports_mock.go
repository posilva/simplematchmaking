// Code generated by MockGen. DO NOT EDIT.
// Source: internal/core/ports/ports.go
//
// Generated by this command:
//
//	mockgen -source=internal/core/ports/ports.go -destination=internal/core/ports/mocks/ports_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	domain "github.com/posilva/simplematchmaking/internal/core/domain"
	ports "github.com/posilva/simplematchmaking/internal/core/ports"
	gomock "go.uber.org/mock/gomock"
)

// MockCodec is a mock of Codec interface.
type MockCodec struct {
	ctrl     *gomock.Controller
	recorder *MockCodecMockRecorder
}

// MockCodecMockRecorder is the mock recorder for MockCodec.
type MockCodecMockRecorder struct {
	mock *MockCodec
}

// NewMockCodec creates a new mock instance.
func NewMockCodec(ctrl *gomock.Controller) *MockCodec {
	mock := &MockCodec{ctrl: ctrl}
	mock.recorder = &MockCodecMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCodec) EXPECT() *MockCodecMockRecorder {
	return m.recorder
}

// Decode mocks base method.
func (m *MockCodec) Decode(data []byte, v any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Decode", data, v)
	ret0, _ := ret[0].(error)
	return ret0
}

// Decode indicates an expected call of Decode.
func (mr *MockCodecMockRecorder) Decode(data, v any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Decode", reflect.TypeOf((*MockCodec)(nil).Decode), data, v)
}

// Encode mocks base method.
func (m *MockCodec) Encode(v any) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Encode", v)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Encode indicates an expected call of Encode.
func (mr *MockCodecMockRecorder) Encode(v any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Encode", reflect.TypeOf((*MockCodec)(nil).Encode), v)
}

// MockLock is a mock of Lock interface.
type MockLock struct {
	ctrl     *gomock.Controller
	recorder *MockLockMockRecorder
}

// MockLockMockRecorder is the mock recorder for MockLock.
type MockLockMockRecorder struct {
	mock *MockLock
}

// NewMockLock creates a new mock instance.
func NewMockLock(ctrl *gomock.Controller) *MockLock {
	mock := &MockLock{ctrl: ctrl}
	mock.recorder = &MockLockMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLock) EXPECT() *MockLockMockRecorder {
	return m.recorder
}

// Acquire mocks base method.
func (m *MockLock) Acquire(ctx context.Context, key string) (context.Context, context.CancelFunc, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Acquire", ctx, key)
	ret0, _ := ret[0].(context.Context)
	ret1, _ := ret[1].(context.CancelFunc)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Acquire indicates an expected call of Acquire.
func (mr *MockLockMockRecorder) Acquire(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Acquire", reflect.TypeOf((*MockLock)(nil).Acquire), ctx, key)
}

// MockMatchResultsListHandler is a mock of MatchResultsListHandler interface.
type MockMatchResultsListHandler struct {
	ctrl     *gomock.Controller
	recorder *MockMatchResultsListHandlerMockRecorder
}

// MockMatchResultsListHandlerMockRecorder is the mock recorder for MockMatchResultsListHandler.
type MockMatchResultsListHandlerMockRecorder struct {
	mock *MockMatchResultsListHandler
}

// NewMockMatchResultsListHandler creates a new mock instance.
func NewMockMatchResultsListHandler(ctrl *gomock.Controller) *MockMatchResultsListHandler {
	mock := &MockMatchResultsListHandler{ctrl: ctrl}
	mock.recorder = &MockMatchResultsListHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMatchResultsListHandler) EXPECT() *MockMatchResultsListHandlerMockRecorder {
	return m.recorder
}

// HandleMatchResultsError mocks base method.
func (m *MockMatchResultsListHandler) HandleMatchResultsError(err error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleMatchResultsError", err)
}

// HandleMatchResultsError indicates an expected call of HandleMatchResultsError.
func (mr *MockMatchResultsListHandlerMockRecorder) HandleMatchResultsError(err any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleMatchResultsError", reflect.TypeOf((*MockMatchResultsListHandler)(nil).HandleMatchResultsError), err)
}

// HandleMatchResultsOK mocks base method.
func (m *MockMatchResultsListHandler) HandleMatchResultsOK(matches []domain.MatchResult) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleMatchResultsOK", matches)
}

// HandleMatchResultsOK indicates an expected call of HandleMatchResultsOK.
func (mr *MockMatchResultsListHandlerMockRecorder) HandleMatchResultsOK(matches any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleMatchResultsOK", reflect.TypeOf((*MockMatchResultsListHandler)(nil).HandleMatchResultsOK), matches)
}

// MockMatchmaker is a mock of Matchmaker interface.
type MockMatchmaker struct {
	ctrl     *gomock.Controller
	recorder *MockMatchmakerMockRecorder
}

// MockMatchmakerMockRecorder is the mock recorder for MockMatchmaker.
type MockMatchmakerMockRecorder struct {
	mock *MockMatchmaker
}

// NewMockMatchmaker creates a new mock instance.
func NewMockMatchmaker(ctrl *gomock.Controller) *MockMatchmaker {
	mock := &MockMatchmaker{ctrl: ctrl}
	mock.recorder = &MockMatchmakerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMatchmaker) EXPECT() *MockMatchmakerMockRecorder {
	return m.recorder
}

// AddPlayer mocks base method.
func (m *MockMatchmaker) AddPlayer(ctx context.Context, ticketID string, p domain.Player) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPlayer", ctx, ticketID, p)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddPlayer indicates an expected call of AddPlayer.
func (mr *MockMatchmakerMockRecorder) AddPlayer(ctx, ticketID, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPlayer", reflect.TypeOf((*MockMatchmaker)(nil).AddPlayer), ctx, ticketID, p)
}

// Matchmake mocks base method.
func (m *MockMatchmaker) Matchmake() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Matchmake")
}

// Matchmake indicates an expected call of Matchmake.
func (mr *MockMatchmakerMockRecorder) Matchmake() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Matchmake", reflect.TypeOf((*MockMatchmaker)(nil).Matchmake))
}

// Subscribe mocks base method.
func (m *MockMatchmaker) Subscribe(handler ports.MatchResultsListHandler) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Subscribe", handler)
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockMatchmakerMockRecorder) Subscribe(handler any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockMatchmaker)(nil).Subscribe), handler)
}

// MockQueue is a mock of Queue interface.
type MockQueue struct {
	ctrl     *gomock.Controller
	recorder *MockQueueMockRecorder
}

// MockQueueMockRecorder is the mock recorder for MockQueue.
type MockQueueMockRecorder struct {
	mock *MockQueue
}

// NewMockQueue creates a new mock instance.
func NewMockQueue(ctrl *gomock.Controller) *MockQueue {
	mock := &MockQueue{ctrl: ctrl}
	mock.recorder = &MockQueueMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQueue) EXPECT() *MockQueueMockRecorder {
	return m.recorder
}

// Enqueue mocks base method.
func (m *MockQueue) Enqueue(ctx context.Context, qe domain.QueueEntry) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Enqueue", ctx, qe)
	ret0, _ := ret[0].(error)
	return ret0
}

// Enqueue indicates an expected call of Enqueue.
func (mr *MockQueueMockRecorder) Enqueue(ctx, qe any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Enqueue", reflect.TypeOf((*MockQueue)(nil).Enqueue), ctx, qe)
}

// Make mocks base method.
func (m *MockQueue) Make(ctx context.Context) ([]domain.MatchResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Make", ctx)
	ret0, _ := ret[0].([]domain.MatchResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Make indicates an expected call of Make.
func (mr *MockQueueMockRecorder) Make(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Make", reflect.TypeOf((*MockQueue)(nil).Make), ctx)
}

// Name mocks base method.
func (m *MockQueue) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockQueueMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockQueue)(nil).Name))
}

// MockMatchmakingService is a mock of MatchmakingService interface.
type MockMatchmakingService struct {
	ctrl     *gomock.Controller
	recorder *MockMatchmakingServiceMockRecorder
}

// MockMatchmakingServiceMockRecorder is the mock recorder for MockMatchmakingService.
type MockMatchmakingServiceMockRecorder struct {
	mock *MockMatchmakingService
}

// NewMockMatchmakingService creates a new mock instance.
func NewMockMatchmakingService(ctrl *gomock.Controller) *MockMatchmakingService {
	mock := &MockMatchmakingService{ctrl: ctrl}
	mock.recorder = &MockMatchmakingServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMatchmakingService) EXPECT() *MockMatchmakingServiceMockRecorder {
	return m.recorder
}

// CancelMatch mocks base method.
func (m *MockMatchmakingService) CancelMatch(ctx context.Context, ticketID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CancelMatch", ctx, ticketID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CancelMatch indicates an expected call of CancelMatch.
func (mr *MockMatchmakingServiceMockRecorder) CancelMatch(ctx, ticketID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelMatch", reflect.TypeOf((*MockMatchmakingService)(nil).CancelMatch), ctx, ticketID)
}

// CheckMatch mocks base method.
func (m *MockMatchmakingService) CheckMatch(ctx context.Context, ticketID string) (domain.Match, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckMatch", ctx, ticketID)
	ret0, _ := ret[0].(domain.Match)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckMatch indicates an expected call of CheckMatch.
func (mr *MockMatchmakingServiceMockRecorder) CheckMatch(ctx, ticketID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckMatch", reflect.TypeOf((*MockMatchmakingService)(nil).CheckMatch), ctx, ticketID)
}

// FindMatch mocks base method.
func (m *MockMatchmakingService) FindMatch(ctx context.Context, queue string, p domain.Player) (domain.Ticket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindMatch", ctx, queue, p)
	ret0, _ := ret[0].(domain.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindMatch indicates an expected call of FindMatch.
func (mr *MockMatchmakingServiceMockRecorder) FindMatch(ctx, queue, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindMatch", reflect.TypeOf((*MockMatchmakingService)(nil).FindMatch), ctx, queue, p)
}

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// DeletePlayerSlot mocks base method.
func (m *MockRepository) DeletePlayerSlot(ctx context.Context, playerID, slot string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePlayerSlot", ctx, playerID, slot)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePlayerSlot indicates an expected call of DeletePlayerSlot.
func (mr *MockRepositoryMockRecorder) DeletePlayerSlot(ctx, playerID, slot any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePlayerSlot", reflect.TypeOf((*MockRepository)(nil).DeletePlayerSlot), ctx, playerID, slot)
}

// DeleteTicket mocks base method.
func (m *MockRepository) DeleteTicket(ctx context.Context, ticketID string) (domain.TicketRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTicket", ctx, ticketID)
	ret0, _ := ret[0].(domain.TicketRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteTicket indicates an expected call of DeleteTicket.
func (mr *MockRepositoryMockRecorder) DeleteTicket(ctx, ticketID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTicket", reflect.TypeOf((*MockRepository)(nil).DeleteTicket), ctx, ticketID)
}

// GetTicket mocks base method.
func (m *MockRepository) GetTicket(ctx context.Context, ticketID string) (domain.TicketRecord, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTicket", ctx, ticketID)
	ret0, _ := ret[0].(domain.TicketRecord)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTicket indicates an expected call of GetTicket.
func (mr *MockRepositoryMockRecorder) GetTicket(ctx, ticketID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTicket", reflect.TypeOf((*MockRepository)(nil).GetTicket), ctx, ticketID)
}

// ReservePlayerSlot mocks base method.
func (m *MockRepository) ReservePlayerSlot(ctx context.Context, playerID, slot, ticketID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReservePlayerSlot", ctx, playerID, slot, ticketID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReservePlayerSlot indicates an expected call of ReservePlayerSlot.
func (mr *MockRepositoryMockRecorder) ReservePlayerSlot(ctx, playerID, slot, ticketID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReservePlayerSlot", reflect.TypeOf((*MockRepository)(nil).ReservePlayerSlot), ctx, playerID, slot, ticketID)
}

// UpdateTicket mocks base method.
func (m *MockRepository) UpdateTicket(ctx context.Context, ticket domain.TicketRecord) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTicket", ctx, ticket)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTicket indicates an expected call of UpdateTicket.
func (mr *MockRepositoryMockRecorder) UpdateTicket(ctx, ticket any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTicket", reflect.TypeOf((*MockRepository)(nil).UpdateTicket), ctx, ticket)
}

// MockLogger is a mock of Logger interface.
type MockLogger struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerMockRecorder
}

// MockLoggerMockRecorder is the mock recorder for MockLogger.
type MockLoggerMockRecorder struct {
	mock *MockLogger
}

// NewMockLogger creates a new mock instance.
func NewMockLogger(ctrl *gomock.Controller) *MockLogger {
	mock := &MockLogger{ctrl: ctrl}
	mock.recorder = &MockLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogger) EXPECT() *MockLoggerMockRecorder {
	return m.recorder
}

// Debug mocks base method.
func (m *MockLogger) Debug(msg string, v ...any) error {
	m.ctrl.T.Helper()
	varargs := []any{msg}
	for _, a := range v {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Debug", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Debug indicates an expected call of Debug.
func (mr *MockLoggerMockRecorder) Debug(msg any, v ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{msg}, v...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debug", reflect.TypeOf((*MockLogger)(nil).Debug), varargs...)
}

// Error mocks base method.
func (m *MockLogger) Error(msg string, v ...any) error {
	m.ctrl.T.Helper()
	varargs := []any{msg}
	for _, a := range v {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Error", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Error indicates an expected call of Error.
func (mr *MockLoggerMockRecorder) Error(msg any, v ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{msg}, v...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockLogger)(nil).Error), varargs...)
}

// Info mocks base method.
func (m *MockLogger) Info(msg string, v ...any) error {
	m.ctrl.T.Helper()
	varargs := []any{msg}
	for _, a := range v {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Info", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Info indicates an expected call of Info.
func (mr *MockLoggerMockRecorder) Info(msg any, v ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{msg}, v...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockLogger)(nil).Info), varargs...)
}

// MockProvider is a mock of Provider interface.
type MockProvider[T any] struct {
	ctrl     *gomock.Controller
	recorder *MockProviderMockRecorder[T]
}

// MockProviderMockRecorder is the mock recorder for MockProvider.
type MockProviderMockRecorder[T any] struct {
	mock *MockProvider[T]
}

// NewMockProvider creates a new mock instance.
func NewMockProvider[T any](ctrl *gomock.Controller) *MockProvider[T] {
	mock := &MockProvider[T]{ctrl: ctrl}
	mock.recorder = &MockProviderMockRecorder[T]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProvider[T]) EXPECT() *MockProviderMockRecorder[T] {
	return m.recorder
}

// Provide mocks base method.
func (m *MockProvider[T]) Provide() (T, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Provide")
	ret0, _ := ret[0].(T)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Provide indicates an expected call of Provide.
func (mr *MockProviderMockRecorder[T]) Provide() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Provide", reflect.TypeOf((*MockProvider[T])(nil).Provide))
}

// MockTelemetryReporter is a mock of TelemetryReporter interface.
type MockTelemetryReporter struct {
	ctrl     *gomock.Controller
	recorder *MockTelemetryReporterMockRecorder
}

// MockTelemetryReporterMockRecorder is the mock recorder for MockTelemetryReporter.
type MockTelemetryReporterMockRecorder struct {
	mock *MockTelemetryReporter
}

// NewMockTelemetryReporter creates a new mock instance.
func NewMockTelemetryReporter(ctrl *gomock.Controller) *MockTelemetryReporter {
	mock := &MockTelemetryReporter{ctrl: ctrl}
	mock.recorder = &MockTelemetryReporterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTelemetryReporter) EXPECT() *MockTelemetryReporterMockRecorder {
	return m.recorder
}

// ReportCounter mocks base method.
func (m *MockTelemetryReporter) ReportCounter(name string, value float64, tags map[string]string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ReportCounter", name, value, tags)
}

// ReportCounter indicates an expected call of ReportCounter.
func (mr *MockTelemetryReporterMockRecorder) ReportCounter(name, value, tags any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReportCounter", reflect.TypeOf((*MockTelemetryReporter)(nil).ReportCounter), name, value, tags)
}

// ReportGauge mocks base method.
func (m *MockTelemetryReporter) ReportGauge(name string, value float64, tags map[string]string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ReportGauge", name, value, tags)
}

// ReportGauge indicates an expected call of ReportGauge.
func (mr *MockTelemetryReporterMockRecorder) ReportGauge(name, value, tags any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReportGauge", reflect.TypeOf((*MockTelemetryReporter)(nil).ReportGauge), name, value, tags)
}

// ReportHistogram mocks base method.
func (m *MockTelemetryReporter) ReportHistogram(name string, value float64, tags map[string]string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ReportHistogram", name, value, tags)
}

// ReportHistogram indicates an expected call of ReportHistogram.
func (mr *MockTelemetryReporterMockRecorder) ReportHistogram(name, value, tags any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReportHistogram", reflect.TypeOf((*MockTelemetryReporter)(nil).ReportHistogram), name, value, tags)
}

// ReportSummary mocks base method.
func (m *MockTelemetryReporter) ReportSummary(name string, value float64, tags map[string]string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ReportSummary", name, value, tags)
}

// ReportSummary indicates an expected call of ReportSummary.
func (mr *MockTelemetryReporterMockRecorder) ReportSummary(name, value, tags any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReportSummary", reflect.TypeOf((*MockTelemetryReporter)(nil).ReportSummary), name, value, tags)
}

// SetDefaultTags mocks base method.
func (m *MockTelemetryReporter) SetDefaultTags(tags map[string]string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetDefaultTags", tags)
}

// SetDefaultTags indicates an expected call of SetDefaultTags.
func (mr *MockTelemetryReporterMockRecorder) SetDefaultTags(tags any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDefaultTags", reflect.TypeOf((*MockTelemetryReporter)(nil).SetDefaultTags), tags)
}
