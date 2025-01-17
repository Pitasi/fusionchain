// Copyright (c) Fusion Laboratories LTD
// SPDX-License-Identifier: BUSL-1.1
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMsgAppendChildWorkspace_NewMsgAppendChildWorkspace(t *testing.T) {

	tests := []struct {
		name string
		msg  *MsgAppendChildWorkspace
		err  error
	}{
		{
			name: "PASS: happy path",
			msg: &MsgAppendChildWorkspace{
				Creator:             sample.AccAddress(),
				ParentWorkspaceAddr: "qredoworkspace14a2hpadpsy9h5m6us54",
				ChildWorkspaceAddr:  "qredoworkspace10j06zdk5gyl6vrss5d5",
				Btl:                 100,
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMsgAppendChildWorkspace(tt.msg.Creator, tt.msg.ParentWorkspaceAddr, tt.msg.ChildWorkspaceAddr, 100)

			assert.Equalf(t, tt.msg, got, "want", tt.msg)
		})
	}
}

func TestMsgAppendChildWorkspace_Route(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgAppendChildWorkspace
	}{
		{
			name: "PASS: happy path",
			msg: MsgAppendChildWorkspace{
				Creator:             sample.AccAddress(),
				ParentWorkspaceAddr: "qredoworkspace14a2hpadpsy9h5m6us54",
				ChildWorkspaceAddr:  "qredoworkspace10j06zdk5gyl6vrss5d5",
				Btl:                 100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, ModuleName, tt.msg.Route(), "Route()")
		})
	}
}

func TestMsgAppendChildWorkspace_Type(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgAppendChildWorkspace
	}{
		{
			name: "PASS: valid address",
			msg: MsgAppendChildWorkspace{
				Creator:             sample.AccAddress(),
				ParentWorkspaceAddr: "qredoworkspace14a2hpadpsy9h5m6us54",
				ChildWorkspaceAddr:  "qredoworkspace10j06zdk5gyl6vrss5d5",
				Btl:                 100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, TypeMsgAppendChildWorkspace, tt.msg.Type(), "Type()")
		})
	}
}

func TestMsgAppendChildWorkspace_GetSigners(t *testing.T) {

	tests := []struct {
		name string
		msg  *MsgAppendChildWorkspace
	}{
		{
			name: "PASS: happy path",
			msg: &MsgAppendChildWorkspace{
				Creator:             "qredo1n7x7nv2urvdtc36tvhvc4dg6wfnnwh3cmt9j9w",
				ParentWorkspaceAddr: "qredoworkspace14a2hpadpsy9h5m6us54",
				ChildWorkspaceAddr:  "qredoworkspace10j06zdk5gyl6vrss5d5",
				Btl:                 100,
			},
		},
		{
			name: "FAIL: invalid signer",
			msg: &MsgAppendChildWorkspace{
				Creator:             "invalid",
				ParentWorkspaceAddr: "qredoworkspace14a2hpadpsy9h5m6us54",
				ChildWorkspaceAddr:  "qredoworkspace10j06zdk5gyl6vrss5d5",
				Btl:                 100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc, err := sdk.AccAddressFromBech32(tt.msg.Creator)
			if err != nil {
				assert.Panics(t, func() { tt.msg.GetSigners() })
			} else {
				msg := NewMsgAppendChildWorkspace(tt.msg.Creator, tt.msg.ParentWorkspaceAddr, tt.msg.ChildWorkspaceAddr, tt.msg.Btl)
				got := msg.GetSigners()

				assert.Equal(t, []sdk.AccAddress{acc}, got)
			}
		})
	}
}

func TestMsgAppendChildWorkspace_GetSignBytes(t *testing.T) {

	tests := []struct {
		name string
		msg  *MsgAppendChildWorkspace
	}{
		{
			name: "PASS: happy path",
			msg: &MsgAppendChildWorkspace{
				Creator:             "qredo1nexzt4fcc84mgnqwjdhxg6veu97eyy9rgzkczs",
				ParentWorkspaceAddr: "qredoworkspace14a2hpadpsy9h5m6us54",
				ChildWorkspaceAddr:  "qredoworkspace10j06zdk5gyl6vrss5d5",
				Btl:                 100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := NewMsgAppendChildWorkspace(tt.msg.Creator, tt.msg.ParentWorkspaceAddr, tt.msg.ChildWorkspaceAddr, tt.msg.Btl)
			got := msg.GetSignBytes()

			bz := ModuleCdc.MustMarshalJSON(msg)
			sortedBz := sdk.MustSortJSON(bz)

			require.Equal(t, sortedBz, got, "GetSignBytes() result doesn't match sorted JSON bytes")

		})
	}
}

func TestMsgAppendChildWorkspace_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgAppendChildWorkspace
		err  error
	}{
		{
			name: "FAIL: invalid address",
			msg: MsgAppendChildWorkspace{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "PASS: valid address",
			msg: MsgAppendChildWorkspace{
				Creator: sample.AccAddress(),
			},
		},
		// todo: add tests for workspace addresses. requires a new function that checks if workspace address is valid
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
