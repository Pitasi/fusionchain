// Copyright (c) Fusion Laboratories LTD
// SPDX-License-Identifier: BUSL-1.1
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMsgNewWorkspace_NewMsgNewWorkspace(t *testing.T) {

	tests := []struct {
		name string
		msg  *MsgNewWorkspace
	}{
		{
			name: "PASS: happy path",
			msg: &MsgNewWorkspace{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMsgNewWorkspace(tt.msg.Creator, tt.msg.AdminPolicyId, tt.msg.SignPolicyId)

			assert.Equalf(t, tt.msg, got, "want", tt.msg)
		})
	}
}

func TestMsgNewWorkspace_Route(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgNewWorkspace
	}{
		{
			name: "PASS: valid address",
			msg: MsgNewWorkspace{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, ModuleName, tt.msg.Route(), "Route()")
		})
	}
}

func TestMsgNewWorkspace_Type(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgNewWorkspace
	}{
		{
			name: "PASS: valid address",
			msg: MsgNewWorkspace{
				Creator: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, TypeMsgNewWorkspace, tt.msg.Type(), "Type()")
		})
	}
}

func TestMsgNewWorkspace_GetSigners(t *testing.T) {

	tests := []struct {
		name string
		msg  *MsgNewWorkspace
	}{
		{
			name: "PASS: happy path",
			msg: &MsgNewWorkspace{
				Creator: "qredo1n7x7nv2urvdtc36tvhvc4dg6wfnnwh3cmt9j9w",
			},
		},
		{
			name: "FAIL: invalid signer",
			msg: &MsgNewWorkspace{
				Creator: "invalid",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc, err := sdk.AccAddressFromBech32(tt.msg.Creator)
			if err != nil {
				assert.Panics(t, func() { tt.msg.GetSigners() })
			} else {
				msg := NewMsgNewWorkspace(tt.msg.Creator, tt.msg.AdminPolicyId, tt.msg.SignPolicyId)
				got := msg.GetSigners()

				assert.Equal(t, []sdk.AccAddress{acc}, got)
			}
		})
	}
}

func TestMsgNewWorkspace_GetSignBytes(t *testing.T) {

	tests := []struct {
		name string
		msg  *MsgNewWorkspace
	}{
		{
			name: "PASS: happy path",
			msg: &MsgNewWorkspace{
				Creator:       "qredo1nexzt4fcc84mgnqwjdhxg6veu97eyy9rgzkczs",
				AdminPolicyId: 0,
				SignPolicyId:  0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := NewMsgNewWorkspace(tt.msg.Creator, tt.msg.AdminPolicyId, tt.msg.SignPolicyId)
			got := msg.GetSignBytes()

			bz := ModuleCdc.MustMarshalJSON(msg)
			sortedBz := sdk.MustSortJSON(bz)

			require.Equal(t, sortedBz, got, "GetSignBytes() result doesn't match sorted JSON bytes")

		})
	}
}

func TestMsgNewWorkspace_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgNewWorkspace
		err  error
	}{
		{
			name: "FAIL: invalid address",
			msg: MsgNewWorkspace{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "PASS: valid address",
			msg: MsgNewWorkspace{
				Creator: sample.AccAddress(),
			},
		},
		{
			name: "PASS: additional owners",
			msg: MsgNewWorkspace{
				Creator: sample.AccAddress(),
				AdditionalOwners: []string{
					sample.AccAddress(),
				},
			},
		},
		{
			name: "FAIL: invalid additional owners",
			msg: MsgNewWorkspace{
				Creator: sample.AccAddress(),
				AdditionalOwners: []string{
					"invalid_Owner",
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "FAIL: duplicated additional owners",
			msg: MsgNewWorkspace{
				Creator: "qredo1n7x7nv2urvdtc36tvhvc4dg6wfnnwh3cmt9j9w",
				AdditionalOwners: []string{
					"qredo1n7x7nv2urvdtc36tvhvc4dg6wfnnwh3cmt9j9w",
				},
			},
			err: ErrDuplicateOwners,
		},
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
