package types_test

import (
	"fmt"
	. "github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

var exampleMetadataV13 = Metadata{
	MagicNumber:   0x6174656d,
	Version:       13,
	IsMetadataV13: true,
	AsMetadataV13: exampleRuntimeMetadataV13,
}

var exampleRuntimeMetadataV13 = MetadataV13{
	Modules: []ModuleMetadataV13{exampleModuleMetadataV13Empty, exampleModuleMetadataV131, exampleModuleMetadataV132},
}

var exampleModuleMetadataV13Empty = ModuleMetadataV13{
	Name:       "EmptyModule",
	HasStorage: false,
	Storage:    StorageMetadataV13{},
	HasCalls:   false,
	Calls:      nil,
	HasEvents:  false,
	Events:     nil,
	Constants:  nil,
	Errors:     nil,
	Index:      0,
}

var exampleModuleMetadataV131 = ModuleMetadataV13{
	Name:       "Module1",
	HasStorage: true,
	Storage:    exampleStorageMetadataV13,
	HasCalls:   true,
	Calls:      []FunctionMetadataV4{exampleFunctionMetadataV4},
	HasEvents:  true,
	Events:     []EventMetadataV4{exampleEventMetadataV4},
	Constants:  []ModuleConstantMetadataV6{exampleModuleConstantMetadataV6},
	Errors:     []ErrorMetadataV8{exampleErrorMetadataV8},
	Index:      1,
}

var exampleModuleMetadataV132 = ModuleMetadataV13{
	Name:       "Module2",
	HasStorage: true,
	Storage:    exampleStorageMetadataV13,
	HasCalls:   true,
	Calls:      []FunctionMetadataV4{exampleFunctionMetadataV4},
	HasEvents:  true,
	Events:     []EventMetadataV4{exampleEventMetadataV4},
	Constants:  []ModuleConstantMetadataV6{exampleModuleConstantMetadataV6},
	Errors:     []ErrorMetadataV8{exampleErrorMetadataV8},
	Index:      2,
}

var exampleStorageMetadataV13 = StorageMetadataV13{
	Prefix: "myStoragePrefix",
	Items: []StorageEntryMetadataV13{exampleStorageFunctionMetadataV13Type, exampleStorageFunctionMetadataV13Map,
		exampleStorageFunctionMetadataV13DoubleMap, exampleStorageFunctionMetadataV13NMap},
}

var exampleStorageFunctionMetadataV13Type = StorageEntryMetadataV13{
	Name:          "myStorageFunc",
	Modifier:      StorageFunctionModifierV0{IsOptional: true},
	Type:          StorageEntryTypeV13{IsPlain: true, AsPlain: "U8"},
	Fallback:      []byte{23, 14},
	Documentation: []Text{"My", "storage func", "doc"},
}

var exampleStorageFunctionMetadataV13Map = StorageEntryMetadataV13{
	Name:          "myStorageFunc2",
	Modifier:      StorageFunctionModifierV0{IsOptional: true},
	Type:          StorageEntryTypeV13{IsMap: true, AsMap: exampleMapTypeV13},
	Fallback:      []byte{23, 14},
	Documentation: []Text{"My", "storage func", "doc"},
}

var exampleStorageFunctionMetadataV13DoubleMap = StorageEntryMetadataV13{
	Name:          "myStorageFunc3",
	Modifier:      StorageFunctionModifierV0{IsOptional: true},
	Type:          StorageEntryTypeV13{IsDoubleMap: true, AsDoubleMap: exampleDoubleMapTypeV13},
	Fallback:      []byte{23, 14},
	Documentation: []Text{"My", "storage func", "doc"},
}

var exampleStorageFunctionMetadataV13NMap = StorageEntryMetadataV13{
	Name:          "myStorageFunc4",
	Modifier:      StorageFunctionModifierV0{IsOptional: true},
	Type:          StorageEntryTypeV13{IsNMap: true, AsNMap: exampleNMapTypeV13},
	Fallback:      []byte{23, 14},
	Documentation: []Text{"My", "storage func", "doc"},
}

var exampleMapTypeV13 = MapTypeV10{
	Hasher: StorageHasherV10{IsBlake2_256: true},
	Key:    "my key",
	Value:  "and my value",
	Linked: false,
}

var exampleDoubleMapTypeV13 = DoubleMapTypeV10{
	Hasher:     StorageHasherV10{IsBlake2_256: true},
	Key1:       "myKey",
	Key2:       "otherKey",
	Value:      "and a value",
	Key2Hasher: StorageHasherV10{IsTwox256: true},
}

var exampleNMapTypeV13 = NMapTypeV13 {
	KeyVec: []Type{"AssetId", "AccountId", "AccountId"},
	Hashers: []StorageHasherV10{{IsBlake2_128: true}},
	Value: "AssetApproval",
}

func TestNewMetadataV13_Decode(t *testing.T) {
	tests := []struct {
		name, hexData string
	}{
		{
			"PolkadotV13", ExamplaryMetadataV13PolkadotString,
		},
	}

	for _, s := range tests {
		t.Run(s.name, func(t *testing.T) {
			metadata := NewMetadataV13()
			err := DecodeFromBytes(MustHexDecodeString(s.hexData), metadata)
			assert.True(t, metadata.IsMetadataV13)
			assert.NoError(t, err)
			data, err := EncodeToBytes(metadata)
			assert.NoError(t, err)
			assert.Equal(t, s.hexData, HexEncodeToString(data))
		})
	}
}

func TestMetadataV13_ExistsModuleMetadata(t *testing.T) {
	assert.True(t, exampleMetadataV13.ExistsModuleMetadata("EmptyModule"))
	assert.False(t, exampleMetadataV13.ExistsModuleMetadata("NotExistModule"))
}

func TestMetadataV13_TestFindCallIndex(t *testing.T) {
	callIndex, err := exampleMetadataV13.FindCallIndex("Module2.my function")
	assert.NoError(t, err)
	assert.Equal(t, exampleModuleMetadataV132.Index, callIndex.SectionIndex)
	assert.Equal(t, uint8(0), callIndex.MethodIndex)
}

func TestMetadataV13_FindEventNamesForEventID(t *testing.T) {
	module, event, err := exampleMetadataV13.FindEventNamesForEventID(EventID([2]byte{1, 0}))

	assert.NoError(t, err)
	assert.Equal(t, exampleModuleMetadataV131.Name, module)
	assert.Equal(t, exampleEventMetadataV4.Name, event)
}

func TestMetadataV13_TestFindStorageEntryMetadata(t *testing.T) {
	storageMap, err := exampleMetadataV13.FindStorageEntryMetadata("myStoragePrefix", "myStorageFunc2")
	fmt.Println("storageMap:", storageMap)
	assert.NoError(t, err)
}
