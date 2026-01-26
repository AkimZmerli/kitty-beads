// Package issues provides the issue management vertical slice.
package issues

import (
	"context"
	"encoding/json"

	"github.com/steveyegge/beads/internal/utils"
)

// RPCHandler provides RPC handlers for issue operations.
// It wraps the Service and handles JSON marshaling/unmarshaling.
type RPCHandler struct {
	service *Service
}

// NewRPCHandler creates a new RPC handler.
func NewRPCHandler(service *Service) *RPCHandler {
	return &RPCHandler{service: service}
}

// RPCResponse represents a standard RPC response.
type RPCResponse struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

func successResponse(data interface{}) RPCResponse {
	jsonData, _ := json.Marshal(data)
	return RPCResponse{
		Success: true,
		Data:    jsonData,
	}
}

func errorResponse(err string) RPCResponse {
	return RPCResponse{
		Success: false,
		Error:   err,
	}
}

// HandleCreate handles the create issue RPC operation.
func (h *RPCHandler) HandleCreate(ctx context.Context, args json.RawMessage, actor string) RPCResponse {
	var createArgs CreateArgs
	if err := json.Unmarshal(args, &createArgs); err != nil {
		return errorResponse("invalid create args: " + err.Error())
	}

	issue, err := h.service.Create(ctx, createArgs, actor)
	if err != nil {
		return errorResponse(err.Error())
	}

	return successResponse(issue)
}

// HandleUpdate handles the update issue RPC operation.
func (h *RPCHandler) HandleUpdate(ctx context.Context, args json.RawMessage, actor string) RPCResponse {
	var updateArgs UpdateArgs
	if err := json.Unmarshal(args, &updateArgs); err != nil {
		return errorResponse("invalid update args: " + err.Error())
	}

	issue, err := h.service.Update(ctx, updateArgs, actor)
	if err != nil {
		return errorResponse(err.Error())
	}

	return successResponse(issue)
}

// HandleClose handles the close issue RPC operation.
func (h *RPCHandler) HandleClose(ctx context.Context, args json.RawMessage, actor string) RPCResponse {
	var closeArgs CloseArgs
	if err := json.Unmarshal(args, &closeArgs); err != nil {
		return errorResponse("invalid close args: " + err.Error())
	}

	result, err := h.service.Close(ctx, closeArgs, actor)
	if err != nil {
		return errorResponse(err.Error())
	}

	// For backward compatibility, return just the closed issue if SuggestNext wasn't requested
	if !closeArgs.SuggestNext {
		return successResponse(result.Closed)
	}

	return successResponse(result)
}

// HandleDelete handles the delete issue RPC operation.
func (h *RPCHandler) HandleDelete(ctx context.Context, args json.RawMessage, actor string) RPCResponse {
	var deleteArgs DeleteArgs
	if err := json.Unmarshal(args, &deleteArgs); err != nil {
		return errorResponse("invalid delete args: " + err.Error())
	}

	result, err := h.service.Delete(ctx, deleteArgs, actor)
	if err != nil {
		return errorResponse(err.Error())
	}

	return successResponse(result)
}

// HandleList handles the list issues RPC operation.
func (h *RPCHandler) HandleList(ctx context.Context, args json.RawMessage) RPCResponse {
	var filter ListFilter
	if err := json.Unmarshal(args, &filter); err != nil {
		return errorResponse("invalid list args: " + err.Error())
	}

	issues, err := h.service.List(ctx, filter)
	if err != nil {
		return errorResponse(err.Error())
	}

	return successResponse(issues)
}

// HandleCount handles the count issues RPC operation.
func (h *RPCHandler) HandleCount(ctx context.Context, args json.RawMessage) RPCResponse {
	var filter CountFilter
	if err := json.Unmarshal(args, &filter); err != nil {
		return errorResponse("invalid count args: " + err.Error())
	}

	result, err := h.service.Count(ctx, filter)
	if err != nil {
		return errorResponse(err.Error())
	}

	return successResponse(result)
}

// HandleShow handles the show issue RPC operation.
func (h *RPCHandler) HandleShow(ctx context.Context, args json.RawMessage) RPCResponse {
	var showArgs ShowArgs
	if err := json.Unmarshal(args, &showArgs); err != nil {
		return errorResponse("invalid show args: " + err.Error())
	}

	details, err := h.service.Show(ctx, showArgs.ID)
	if err != nil {
		return errorResponse(err.Error())
	}

	return successResponse(details)
}

// HandleReady handles the ready work RPC operation.
func (h *RPCHandler) HandleReady(ctx context.Context, args json.RawMessage) RPCResponse {
	var filter ReadyFilter
	if err := json.Unmarshal(args, &filter); err != nil {
		return errorResponse("invalid ready args: " + err.Error())
	}

	issues, err := h.service.GetReady(ctx, filter)
	if err != nil {
		return errorResponse(err.Error())
	}

	return successResponse(issues)
}

// HandleBlocked handles the blocked issues RPC operation.
func (h *RPCHandler) HandleBlocked(ctx context.Context, args json.RawMessage) RPCResponse {
	var filter BlockedFilter
	if err := json.Unmarshal(args, &filter); err != nil {
		return errorResponse("invalid blocked args: " + err.Error())
	}

	issues, err := h.service.GetBlocked(ctx, filter)
	if err != nil {
		return errorResponse(err.Error())
	}

	return successResponse(issues)
}

// HandleStale handles the stale issues RPC operation.
func (h *RPCHandler) HandleStale(ctx context.Context, args json.RawMessage) RPCResponse {
	var filter StaleFilter
	if err := json.Unmarshal(args, &filter); err != nil {
		return errorResponse("invalid stale args: " + err.Error())
	}

	issues, err := h.service.GetStale(ctx, filter)
	if err != nil {
		return errorResponse(err.Error())
	}

	return successResponse(issues)
}

// HandleStats handles the statistics RPC operation.
func (h *RPCHandler) HandleStats(ctx context.Context) RPCResponse {
	stats, err := h.service.GetStats(ctx)
	if err != nil {
		return errorResponse(err.Error())
	}

	return successResponse(stats)
}

// HandleEpicStatus handles the epic status RPC operation.
func (h *RPCHandler) HandleEpicStatus(ctx context.Context, args json.RawMessage) RPCResponse {
	var filter EpicStatusFilter
	if err := json.Unmarshal(args, &filter); err != nil {
		return errorResponse("invalid epic status args: " + err.Error())
	}

	epics, err := h.service.GetEpicStatus(ctx, filter)
	if err != nil {
		return errorResponse(err.Error())
	}

	return successResponse(epics)
}

// HandleResolveID handles the resolve ID RPC operation.
func (h *RPCHandler) HandleResolveID(ctx context.Context, args json.RawMessage) RPCResponse {
	var resolveArgs ResolveIDArgs
	if err := json.Unmarshal(args, &resolveArgs); err != nil {
		return errorResponse("invalid resolve_id args: " + err.Error())
	}

	resolvedID, err := utils.ResolvePartialID(ctx, h.service.store, resolveArgs.ID)
	if err != nil {
		return errorResponse("failed to resolve ID: " + err.Error())
	}

	return successResponse(resolvedID)
}
