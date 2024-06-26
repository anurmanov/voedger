/*
 * Copyright (c) 2023-present unTill Pro, Ltd.
 */

package state

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/voedger/voedger/pkg/appdef"
	"github.com/voedger/voedger/pkg/istructs"
)

func TestCommandContextStorage(t *testing.T) {
	cmdResBuilder := istructs.NewNullObjectBuilder()

	wsid := istructs.NullWSID
	arg := istructs.NewNullObject()
	unloggedArg := istructs.NewNullObject()

	wsidFunc := SimpleWSIDFunc(wsid)
	commandPrepareArgs := func() istructs.CommandPrepareArgs {
		return istructs.CommandPrepareArgs{
			PrepareArgs: istructs.PrepareArgs{
				Workpiece:      nil,
				ArgumentObject: arg,
				WSID:           wsid,
				Workspace:      nil,
			},
			ArgumentUnloggedObject: unloggedArg,
		}
	}
	argFunc := func() istructs.IObject { return arg }
	unloggedArgFunc := func() istructs.IObject { return unloggedArg }
	wlogOffsetFunc := func() istructs.Offset { return 42 }

	s := ProvideCommandProcessorStateFactory()(context.Background(), nil, nil, wsidFunc, nil, nil, nil, nil, 1,
		func() istructs.IObjectBuilder { return cmdResBuilder }, commandPrepareArgs, argFunc, unloggedArgFunc, wlogOffsetFunc)

	kb, err := s.KeyBuilder(CommandContext, appdef.NullQName)
	require.NoError(t, err)

	v, err := s.MustExist(kb)
	require.NoError(t, err)
	require.Equal(t, int64(wsid), v.AsInt64(Field_Workspace))
	require.NotNil(t, v.AsValue(Field_ArgumentObject))
	require.NotNil(t, v.AsValue(Field_ArgumentUnloggedObject))
	require.Equal(t, int64(42), v.AsInt64(Field_WLogOffset))

}
