// Copyright (c) Fusion Laboratories LTD
// SPDX-License-Identifier: BUSL-1.1
)

func TestMsgNewPolicy_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgNewPolicy
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgNewPolicy{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgNewPolicy{
				Creator: sample.AccAddress(),
			},
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
