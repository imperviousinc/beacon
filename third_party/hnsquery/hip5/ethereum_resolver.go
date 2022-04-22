// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package hip5

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// DNSResolverABI is the input ABI used to generate the binding from.
const DNSResolverABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"name\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"resource\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"record\",\"type\":\"bytes\"}],\"name\":\"DNSRecordChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"name\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"resource\",\"type\":\"uint16\"}],\"name\":\"DNSRecordDeleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"}],\"name\":\"DNSZoneCleared\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"lastzonehash\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"zonehash\",\"type\":\"bytes\"}],\"name\":\"DNSZonehashChanged\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"}],\"name\":\"clearDNSZone\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"name\",\"type\":\"bytes32\"},{\"internalType\":\"uint16\",\"name\":\"resource\",\"type\":\"uint16\"}],\"name\":\"dnsRecord\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"name\",\"type\":\"bytes32\"}],\"name\":\"hasDNSRecords\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"setDNSRecords\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"hash\",\"type\":\"bytes\"}],\"name\":\"setZonehash\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceID\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"node\",\"type\":\"bytes32\"}],\"name\":\"zonehash\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// DNSResolver is an auto generated Go binding around an Ethereum contract.
type DNSResolver struct {
	DNSResolverCaller     // Read-only binding to the contract
	DNSResolverTransactor // Write-only binding to the contract
	DNSResolverFilterer   // Log filterer for contract events
}

// DNSResolverCaller is an auto generated read-only Go binding around an Ethereum contract.
type DNSResolverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DNSResolverTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DNSResolverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DNSResolverFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DNSResolverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DNSResolverSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DNSResolverSession struct {
	Contract     *DNSResolver      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DNSResolverCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DNSResolverCallerSession struct {
	Contract *DNSResolverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// DNSResolverTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DNSResolverTransactorSession struct {
	Contract     *DNSResolverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// DNSResolverRaw is an auto generated low-level Go binding around an Ethereum contract.
type DNSResolverRaw struct {
	Contract *DNSResolver // Generic contract binding to access the raw methods on
}

// DNSResolverCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DNSResolverCallerRaw struct {
	Contract *DNSResolverCaller // Generic read-only contract binding to access the raw methods on
}

// DNSResolverTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DNSResolverTransactorRaw struct {
	Contract *DNSResolverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDNSResolver creates a new instance of DNSResolver, bound to a specific deployed contract.
func NewDNSResolver(address common.Address, backend bind.ContractBackend) (*DNSResolver, error) {
	contract, err := bindDNSResolver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DNSResolver{DNSResolverCaller: DNSResolverCaller{contract: contract}, DNSResolverTransactor: DNSResolverTransactor{contract: contract}, DNSResolverFilterer: DNSResolverFilterer{contract: contract}}, nil
}

// NewDNSResolverCaller creates a new read-only instance of DNSResolver, bound to a specific deployed contract.
func NewDNSResolverCaller(address common.Address, caller bind.ContractCaller) (*DNSResolverCaller, error) {
	contract, err := bindDNSResolver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DNSResolverCaller{contract: contract}, nil
}

// NewDNSResolverTransactor creates a new write-only instance of DNSResolver, bound to a specific deployed contract.
func NewDNSResolverTransactor(address common.Address, transactor bind.ContractTransactor) (*DNSResolverTransactor, error) {
	contract, err := bindDNSResolver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DNSResolverTransactor{contract: contract}, nil
}

// NewDNSResolverFilterer creates a new log filterer instance of DNSResolver, bound to a specific deployed contract.
func NewDNSResolverFilterer(address common.Address, filterer bind.ContractFilterer) (*DNSResolverFilterer, error) {
	contract, err := bindDNSResolver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DNSResolverFilterer{contract: contract}, nil
}

// bindDNSResolver binds a generic wrapper to an already deployed contract.
func bindDNSResolver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DNSResolverABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DNSResolver *DNSResolverRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DNSResolver.Contract.DNSResolverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DNSResolver *DNSResolverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DNSResolver.Contract.DNSResolverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DNSResolver *DNSResolverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DNSResolver.Contract.DNSResolverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DNSResolver *DNSResolverCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DNSResolver.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DNSResolver *DNSResolverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DNSResolver.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DNSResolver *DNSResolverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DNSResolver.Contract.contract.Transact(opts, method, params...)
}

// DnsRecord is a free data retrieval call binding the contract method 0xa8fa5682.
//
// Solidity: function dnsRecord(bytes32 node, bytes32 name, uint16 resource) view returns(bytes)
func (_DNSResolver *DNSResolverCaller) DnsRecord(opts *bind.CallOpts, node [32]byte, name [32]byte, resource uint16) ([]byte, error) {
	var out []interface{}
	err := _DNSResolver.contract.Call(opts, &out, "dnsRecord", node, name, resource)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// DnsRecord is a free data retrieval call binding the contract method 0xa8fa5682.
//
// Solidity: function dnsRecord(bytes32 node, bytes32 name, uint16 resource) view returns(bytes)
func (_DNSResolver *DNSResolverSession) DnsRecord(node [32]byte, name [32]byte, resource uint16) ([]byte, error) {
	return _DNSResolver.Contract.DnsRecord(&_DNSResolver.CallOpts, node, name, resource)
}

// DnsRecord is a free data retrieval call binding the contract method 0xa8fa5682.
//
// Solidity: function dnsRecord(bytes32 node, bytes32 name, uint16 resource) view returns(bytes)
func (_DNSResolver *DNSResolverCallerSession) DnsRecord(node [32]byte, name [32]byte, resource uint16) ([]byte, error) {
	return _DNSResolver.Contract.DnsRecord(&_DNSResolver.CallOpts, node, name, resource)
}

// HasDNSRecords is a free data retrieval call binding the contract method 0x4cbf6ba4.
//
// Solidity: function hasDNSRecords(bytes32 node, bytes32 name) view returns(bool)
func (_DNSResolver *DNSResolverCaller) HasDNSRecords(opts *bind.CallOpts, node [32]byte, name [32]byte) (bool, error) {
	var out []interface{}
	err := _DNSResolver.contract.Call(opts, &out, "hasDNSRecords", node, name)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasDNSRecords is a free data retrieval call binding the contract method 0x4cbf6ba4.
//
// Solidity: function hasDNSRecords(bytes32 node, bytes32 name) view returns(bool)
func (_DNSResolver *DNSResolverSession) HasDNSRecords(node [32]byte, name [32]byte) (bool, error) {
	return _DNSResolver.Contract.HasDNSRecords(&_DNSResolver.CallOpts, node, name)
}

// HasDNSRecords is a free data retrieval call binding the contract method 0x4cbf6ba4.
//
// Solidity: function hasDNSRecords(bytes32 node, bytes32 name) view returns(bool)
func (_DNSResolver *DNSResolverCallerSession) HasDNSRecords(node [32]byte, name [32]byte) (bool, error) {
	return _DNSResolver.Contract.HasDNSRecords(&_DNSResolver.CallOpts, node, name)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceID) pure returns(bool)
func (_DNSResolver *DNSResolverCaller) SupportsInterface(opts *bind.CallOpts, interfaceID [4]byte) (bool, error) {
	var out []interface{}
	err := _DNSResolver.contract.Call(opts, &out, "supportsInterface", interfaceID)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceID) pure returns(bool)
func (_DNSResolver *DNSResolverSession) SupportsInterface(interfaceID [4]byte) (bool, error) {
	return _DNSResolver.Contract.SupportsInterface(&_DNSResolver.CallOpts, interfaceID)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceID) pure returns(bool)
func (_DNSResolver *DNSResolverCallerSession) SupportsInterface(interfaceID [4]byte) (bool, error) {
	return _DNSResolver.Contract.SupportsInterface(&_DNSResolver.CallOpts, interfaceID)
}

// Zonehash is a free data retrieval call binding the contract method 0x5c98042b.
//
// Solidity: function zonehash(bytes32 node) view returns(bytes)
func (_DNSResolver *DNSResolverCaller) Zonehash(opts *bind.CallOpts, node [32]byte) ([]byte, error) {
	var out []interface{}
	err := _DNSResolver.contract.Call(opts, &out, "zonehash", node)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Zonehash is a free data retrieval call binding the contract method 0x5c98042b.
//
// Solidity: function zonehash(bytes32 node) view returns(bytes)
func (_DNSResolver *DNSResolverSession) Zonehash(node [32]byte) ([]byte, error) {
	return _DNSResolver.Contract.Zonehash(&_DNSResolver.CallOpts, node)
}

// Zonehash is a free data retrieval call binding the contract method 0x5c98042b.
//
// Solidity: function zonehash(bytes32 node) view returns(bytes)
func (_DNSResolver *DNSResolverCallerSession) Zonehash(node [32]byte) ([]byte, error) {
	return _DNSResolver.Contract.Zonehash(&_DNSResolver.CallOpts, node)
}

// ClearDNSZone is a paid mutator transaction binding the contract method 0xad5780af.
//
// Solidity: function clearDNSZone(bytes32 node) returns()
func (_DNSResolver *DNSResolverTransactor) ClearDNSZone(opts *bind.TransactOpts, node [32]byte) (*types.Transaction, error) {
	return _DNSResolver.contract.Transact(opts, "clearDNSZone", node)
}

// ClearDNSZone is a paid mutator transaction binding the contract method 0xad5780af.
//
// Solidity: function clearDNSZone(bytes32 node) returns()
func (_DNSResolver *DNSResolverSession) ClearDNSZone(node [32]byte) (*types.Transaction, error) {
	return _DNSResolver.Contract.ClearDNSZone(&_DNSResolver.TransactOpts, node)
}

// ClearDNSZone is a paid mutator transaction binding the contract method 0xad5780af.
//
// Solidity: function clearDNSZone(bytes32 node) returns()
func (_DNSResolver *DNSResolverTransactorSession) ClearDNSZone(node [32]byte) (*types.Transaction, error) {
	return _DNSResolver.Contract.ClearDNSZone(&_DNSResolver.TransactOpts, node)
}

// SetDNSRecords is a paid mutator transaction binding the contract method 0x0af179d7.
//
// Solidity: function setDNSRecords(bytes32 node, bytes data) returns()
func (_DNSResolver *DNSResolverTransactor) SetDNSRecords(opts *bind.TransactOpts, node [32]byte, data []byte) (*types.Transaction, error) {
	return _DNSResolver.contract.Transact(opts, "setDNSRecords", node, data)
}

// SetDNSRecords is a paid mutator transaction binding the contract method 0x0af179d7.
//
// Solidity: function setDNSRecords(bytes32 node, bytes data) returns()
func (_DNSResolver *DNSResolverSession) SetDNSRecords(node [32]byte, data []byte) (*types.Transaction, error) {
	return _DNSResolver.Contract.SetDNSRecords(&_DNSResolver.TransactOpts, node, data)
}

// SetDNSRecords is a paid mutator transaction binding the contract method 0x0af179d7.
//
// Solidity: function setDNSRecords(bytes32 node, bytes data) returns()
func (_DNSResolver *DNSResolverTransactorSession) SetDNSRecords(node [32]byte, data []byte) (*types.Transaction, error) {
	return _DNSResolver.Contract.SetDNSRecords(&_DNSResolver.TransactOpts, node, data)
}

// SetZonehash is a paid mutator transaction binding the contract method 0xce3decdc.
//
// Solidity: function setZonehash(bytes32 node, bytes hash) returns()
func (_DNSResolver *DNSResolverTransactor) SetZonehash(opts *bind.TransactOpts, node [32]byte, hash []byte) (*types.Transaction, error) {
	return _DNSResolver.contract.Transact(opts, "setZonehash", node, hash)
}

// SetZonehash is a paid mutator transaction binding the contract method 0xce3decdc.
//
// Solidity: function setZonehash(bytes32 node, bytes hash) returns()
func (_DNSResolver *DNSResolverSession) SetZonehash(node [32]byte, hash []byte) (*types.Transaction, error) {
	return _DNSResolver.Contract.SetZonehash(&_DNSResolver.TransactOpts, node, hash)
}

// SetZonehash is a paid mutator transaction binding the contract method 0xce3decdc.
//
// Solidity: function setZonehash(bytes32 node, bytes hash) returns()
func (_DNSResolver *DNSResolverTransactorSession) SetZonehash(node [32]byte, hash []byte) (*types.Transaction, error) {
	return _DNSResolver.Contract.SetZonehash(&_DNSResolver.TransactOpts, node, hash)
}

// DNSResolverDNSRecordChangedIterator is returned from FilterDNSRecordChanged and is used to iterate over the raw logs and unpacked data for DNSRecordChanged events raised by the DNSResolver contract.
type DNSResolverDNSRecordChangedIterator struct {
	Event *DNSResolverDNSRecordChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DNSResolverDNSRecordChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DNSResolverDNSRecordChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DNSResolverDNSRecordChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DNSResolverDNSRecordChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DNSResolverDNSRecordChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DNSResolverDNSRecordChanged represents a DNSRecordChanged event raised by the DNSResolver contract.
type DNSResolverDNSRecordChanged struct {
	Node     [32]byte
	Name     []byte
	Resource uint16
	Record   []byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDNSRecordChanged is a free log retrieval operation binding the contract event 0x52a608b3303a48862d07a73d82fa221318c0027fbbcfb1b2329bface3f19ff2b.
//
// Solidity: event DNSRecordChanged(bytes32 indexed node, bytes name, uint16 resource, bytes record)
func (_DNSResolver *DNSResolverFilterer) FilterDNSRecordChanged(opts *bind.FilterOpts, node [][32]byte) (*DNSResolverDNSRecordChangedIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _DNSResolver.contract.FilterLogs(opts, "DNSRecordChanged", nodeRule)
	if err != nil {
		return nil, err
	}
	return &DNSResolverDNSRecordChangedIterator{contract: _DNSResolver.contract, event: "DNSRecordChanged", logs: logs, sub: sub}, nil
}

// WatchDNSRecordChanged is a free log subscription operation binding the contract event 0x52a608b3303a48862d07a73d82fa221318c0027fbbcfb1b2329bface3f19ff2b.
//
// Solidity: event DNSRecordChanged(bytes32 indexed node, bytes name, uint16 resource, bytes record)
func (_DNSResolver *DNSResolverFilterer) WatchDNSRecordChanged(opts *bind.WatchOpts, sink chan<- *DNSResolverDNSRecordChanged, node [][32]byte) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _DNSResolver.contract.WatchLogs(opts, "DNSRecordChanged", nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DNSResolverDNSRecordChanged)
				if err := _DNSResolver.contract.UnpackLog(event, "DNSRecordChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDNSRecordChanged is a log parse operation binding the contract event 0x52a608b3303a48862d07a73d82fa221318c0027fbbcfb1b2329bface3f19ff2b.
//
// Solidity: event DNSRecordChanged(bytes32 indexed node, bytes name, uint16 resource, bytes record)
func (_DNSResolver *DNSResolverFilterer) ParseDNSRecordChanged(log types.Log) (*DNSResolverDNSRecordChanged, error) {
	event := new(DNSResolverDNSRecordChanged)
	if err := _DNSResolver.contract.UnpackLog(event, "DNSRecordChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DNSResolverDNSRecordDeletedIterator is returned from FilterDNSRecordDeleted and is used to iterate over the raw logs and unpacked data for DNSRecordDeleted events raised by the DNSResolver contract.
type DNSResolverDNSRecordDeletedIterator struct {
	Event *DNSResolverDNSRecordDeleted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DNSResolverDNSRecordDeletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DNSResolverDNSRecordDeleted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DNSResolverDNSRecordDeleted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DNSResolverDNSRecordDeletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DNSResolverDNSRecordDeletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DNSResolverDNSRecordDeleted represents a DNSRecordDeleted event raised by the DNSResolver contract.
type DNSResolverDNSRecordDeleted struct {
	Node     [32]byte
	Name     []byte
	Resource uint16
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDNSRecordDeleted is a free log retrieval operation binding the contract event 0x03528ed0c2a3ebc993b12ce3c16bb382f9c7d88ef7d8a1bf290eaf35955a1207.
//
// Solidity: event DNSRecordDeleted(bytes32 indexed node, bytes name, uint16 resource)
func (_DNSResolver *DNSResolverFilterer) FilterDNSRecordDeleted(opts *bind.FilterOpts, node [][32]byte) (*DNSResolverDNSRecordDeletedIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _DNSResolver.contract.FilterLogs(opts, "DNSRecordDeleted", nodeRule)
	if err != nil {
		return nil, err
	}
	return &DNSResolverDNSRecordDeletedIterator{contract: _DNSResolver.contract, event: "DNSRecordDeleted", logs: logs, sub: sub}, nil
}

// WatchDNSRecordDeleted is a free log subscription operation binding the contract event 0x03528ed0c2a3ebc993b12ce3c16bb382f9c7d88ef7d8a1bf290eaf35955a1207.
//
// Solidity: event DNSRecordDeleted(bytes32 indexed node, bytes name, uint16 resource)
func (_DNSResolver *DNSResolverFilterer) WatchDNSRecordDeleted(opts *bind.WatchOpts, sink chan<- *DNSResolverDNSRecordDeleted, node [][32]byte) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _DNSResolver.contract.WatchLogs(opts, "DNSRecordDeleted", nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DNSResolverDNSRecordDeleted)
				if err := _DNSResolver.contract.UnpackLog(event, "DNSRecordDeleted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDNSRecordDeleted is a log parse operation binding the contract event 0x03528ed0c2a3ebc993b12ce3c16bb382f9c7d88ef7d8a1bf290eaf35955a1207.
//
// Solidity: event DNSRecordDeleted(bytes32 indexed node, bytes name, uint16 resource)
func (_DNSResolver *DNSResolverFilterer) ParseDNSRecordDeleted(log types.Log) (*DNSResolverDNSRecordDeleted, error) {
	event := new(DNSResolverDNSRecordDeleted)
	if err := _DNSResolver.contract.UnpackLog(event, "DNSRecordDeleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DNSResolverDNSZoneClearedIterator is returned from FilterDNSZoneCleared and is used to iterate over the raw logs and unpacked data for DNSZoneCleared events raised by the DNSResolver contract.
type DNSResolverDNSZoneClearedIterator struct {
	Event *DNSResolverDNSZoneCleared // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DNSResolverDNSZoneClearedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DNSResolverDNSZoneCleared)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DNSResolverDNSZoneCleared)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DNSResolverDNSZoneClearedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DNSResolverDNSZoneClearedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DNSResolverDNSZoneCleared represents a DNSZoneCleared event raised by the DNSResolver contract.
type DNSResolverDNSZoneCleared struct {
	Node [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterDNSZoneCleared is a free log retrieval operation binding the contract event 0xb757169b8492ca2f1c6619d9d76ce22803035c3b1d5f6930dffe7b127c1a1983.
//
// Solidity: event DNSZoneCleared(bytes32 indexed node)
func (_DNSResolver *DNSResolverFilterer) FilterDNSZoneCleared(opts *bind.FilterOpts, node [][32]byte) (*DNSResolverDNSZoneClearedIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _DNSResolver.contract.FilterLogs(opts, "DNSZoneCleared", nodeRule)
	if err != nil {
		return nil, err
	}
	return &DNSResolverDNSZoneClearedIterator{contract: _DNSResolver.contract, event: "DNSZoneCleared", logs: logs, sub: sub}, nil
}

// WatchDNSZoneCleared is a free log subscription operation binding the contract event 0xb757169b8492ca2f1c6619d9d76ce22803035c3b1d5f6930dffe7b127c1a1983.
//
// Solidity: event DNSZoneCleared(bytes32 indexed node)
func (_DNSResolver *DNSResolverFilterer) WatchDNSZoneCleared(opts *bind.WatchOpts, sink chan<- *DNSResolverDNSZoneCleared, node [][32]byte) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _DNSResolver.contract.WatchLogs(opts, "DNSZoneCleared", nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DNSResolverDNSZoneCleared)
				if err := _DNSResolver.contract.UnpackLog(event, "DNSZoneCleared", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDNSZoneCleared is a log parse operation binding the contract event 0xb757169b8492ca2f1c6619d9d76ce22803035c3b1d5f6930dffe7b127c1a1983.
//
// Solidity: event DNSZoneCleared(bytes32 indexed node)
func (_DNSResolver *DNSResolverFilterer) ParseDNSZoneCleared(log types.Log) (*DNSResolverDNSZoneCleared, error) {
	event := new(DNSResolverDNSZoneCleared)
	if err := _DNSResolver.contract.UnpackLog(event, "DNSZoneCleared", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DNSResolverDNSZonehashChangedIterator is returned from FilterDNSZonehashChanged and is used to iterate over the raw logs and unpacked data for DNSZonehashChanged events raised by the DNSResolver contract.
type DNSResolverDNSZonehashChangedIterator struct {
	Event *DNSResolverDNSZonehashChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DNSResolverDNSZonehashChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DNSResolverDNSZonehashChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DNSResolverDNSZonehashChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DNSResolverDNSZonehashChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DNSResolverDNSZonehashChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DNSResolverDNSZonehashChanged represents a DNSZonehashChanged event raised by the DNSResolver contract.
type DNSResolverDNSZonehashChanged struct {
	Node         [32]byte
	Lastzonehash []byte
	Zonehash     []byte
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterDNSZonehashChanged is a free log retrieval operation binding the contract event 0x8f15ed4b723ef428f250961da8315675b507046737e19319fc1a4d81bfe87f85.
//
// Solidity: event DNSZonehashChanged(bytes32 indexed node, bytes lastzonehash, bytes zonehash)
func (_DNSResolver *DNSResolverFilterer) FilterDNSZonehashChanged(opts *bind.FilterOpts, node [][32]byte) (*DNSResolverDNSZonehashChangedIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _DNSResolver.contract.FilterLogs(opts, "DNSZonehashChanged", nodeRule)
	if err != nil {
		return nil, err
	}
	return &DNSResolverDNSZonehashChangedIterator{contract: _DNSResolver.contract, event: "DNSZonehashChanged", logs: logs, sub: sub}, nil
}

// WatchDNSZonehashChanged is a free log subscription operation binding the contract event 0x8f15ed4b723ef428f250961da8315675b507046737e19319fc1a4d81bfe87f85.
//
// Solidity: event DNSZonehashChanged(bytes32 indexed node, bytes lastzonehash, bytes zonehash)
func (_DNSResolver *DNSResolverFilterer) WatchDNSZonehashChanged(opts *bind.WatchOpts, sink chan<- *DNSResolverDNSZonehashChanged, node [][32]byte) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _DNSResolver.contract.WatchLogs(opts, "DNSZonehashChanged", nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DNSResolverDNSZonehashChanged)
				if err := _DNSResolver.contract.UnpackLog(event, "DNSZonehashChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDNSZonehashChanged is a log parse operation binding the contract event 0x8f15ed4b723ef428f250961da8315675b507046737e19319fc1a4d81bfe87f85.
//
// Solidity: event DNSZonehashChanged(bytes32 indexed node, bytes lastzonehash, bytes zonehash)
func (_DNSResolver *DNSResolverFilterer) ParseDNSZonehashChanged(log types.Log) (*DNSResolverDNSZonehashChanged, error) {
	event := new(DNSResolverDNSZonehashChanged)
	if err := _DNSResolver.contract.UnpackLog(event, "DNSZonehashChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
