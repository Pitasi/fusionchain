// Copyright (c) Fusion Laboratories LTD
// SPDX-License-Identifier: BUSL-1.1
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/qredo/fusionchain/policy"
	"github.com/qredo/fusionchain/x/identity/types"
	bbird "github.com/qredo/fusionchain/x/policy/keeper"
	bbirdtypes "github.com/qredo/fusionchain/x/policy/types"
)

func (k msgServer) AppendChildWorkspace(goCtx context.Context, msg *types.MsgAppendChildWorkspace) (*types.MsgAppendChildWorkspaceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	parent := k.GetWorkspace(ctx, msg.ParentWorkspaceAddr)
	if parent == nil {
		return nil, fmt.Errorf("workspace not found")
	}

	if !parent.IsOwner(msg.Creator) {
		return nil, fmt.Errorf("creator is not an owner of the workspace")
	}

	if parent.IsChild(msg.ChildWorkspaceAddr) {
		return nil, fmt.Errorf("new child is already a child workspace")
	}

	act, err := k.policyKeeper.AddAction(ctx, msg.Creator, msg, parent.AdminPolicyId, msg.Btl, nil)
	if err != nil {
		return nil, err
	}
	return k.AppendChildWorkspaceActionHandler(ctx, act, &cdctypes.Any{})
}

func (k msgServer) AppendChildWorkspacePolicyGenerator(ctx sdk.Context, msg *types.MsgAppendChildWorkspace) (policy.Policy, error) {
	parent := k.GetWorkspace(ctx, msg.ParentWorkspaceAddr)
	if parent == nil {
		return nil, fmt.Errorf("workspace not found")
	}

	pol := parent.PolicyAppendChild()
	return pol, nil
}

func (k msgServer) AppendChildWorkspaceActionHandler(ctx sdk.Context, act *bbirdtypes.Action, payload *cdctypes.Any) (*types.MsgAppendChildWorkspaceResponse, error) {
	return bbird.TryExecuteAction(
		k.policyKeeper,
		k.cdc,
		ctx,
		act,
		payload,
		func(ctx sdk.Context, msg *types.MsgAppendChildWorkspace) (*types.MsgAppendChildWorkspaceResponse, error) {
			child := k.GetWorkspace(ctx, msg.ChildWorkspaceAddr)
			parent := k.GetWorkspace(ctx, msg.ParentWorkspaceAddr)

			if child == nil || parent == nil {
				return nil, errors.New("one or more invalid workspace addresses provided")
			}

			parent.AddChild(child)
			k.SetWorkspace(ctx, parent)

			return &types.MsgAppendChildWorkspaceResponse{}, nil
		},
	)
}
