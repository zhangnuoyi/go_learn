// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package token

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// SolTokenMetaData contains all meta data concerning the SolToken contract.
var SolTokenMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"initialSupply\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040526040518060400160405280600881526020017f4d7920546f6b656e0000000000000000000000000000000000000000000000008152505f9081610047919061038b565b506040518060400160405280600381526020017f4d544b00000000000000000000000000000000000000000000000000000000008152506001908161008c919061038b565b50601260025f6101000a81548160ff021916908360ff1602179055503480156100b3575f5ffd5b5060405161140838038061140883398181016040528101906100d59190610488565b60025f9054906101000a900460ff1660ff16600a6100f3919061060f565b816100fe9190610659565b60038190555060035460045f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055505061069a565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f60028204905060018216806101c957607f821691505b6020821081036101dc576101db610185565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f6008830261023e7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82610203565b6102488683610203565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f61028c61028761028284610260565b610269565b610260565b9050919050565b5f819050919050565b6102a583610272565b6102b96102b182610293565b84845461020f565b825550505050565b5f5f905090565b6102d06102c1565b6102db81848461029c565b505050565b5b818110156102fe576102f35f826102c8565b6001810190506102e1565b5050565b601f82111561034357610314816101e2565b61031d846101f4565b8101602085101561032c578190505b610340610338856101f4565b8301826102e0565b50505b505050565b5f82821c905092915050565b5f6103635f1984600802610348565b1980831691505092915050565b5f61037b8383610354565b9150826002028217905092915050565b6103948261014e565b67ffffffffffffffff8111156103ad576103ac610158565b5b6103b782546101b2565b6103c2828285610302565b5f60209050601f8311600181146103f3575f84156103e1578287015190505b6103eb8582610370565b865550610452565b601f198416610401866101e2565b5f5b8281101561042857848901518255600182019150602085019450602081019050610403565b868310156104455784890151610441601f891682610354565b8355505b6001600288020188555050505b505050505050565b5f5ffd5b61046781610260565b8114610471575f5ffd5b50565b5f815190506104828161045e565b92915050565b5f6020828403121561049d5761049c61045a565b5b5f6104aa84828501610474565b91505092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f8160011c9050919050565b5f5f8291508390505b600185111561053557808604811115610511576105106104b3565b5b60018516156105205780820291505b808102905061052e856104e0565b94506104f5565b94509492505050565b5f8261054d5760019050610608565b8161055a575f9050610608565b8160018114610570576002811461057a576105a9565b6001915050610608565b60ff84111561058c5761058b6104b3565b5b8360020a9150848211156105a3576105a26104b3565b5b50610608565b5060208310610133831016604e8410600b84101617156105de5782820a9050838111156105d9576105d86104b3565b5b610608565b6105eb84848460016104ec565b92509050818404811115610602576106016104b3565b5b81810290505b9392505050565b5f61061982610260565b915061062483610260565b92506106517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff848461053e565b905092915050565b5f61066382610260565b915061066e83610260565b925082820261067c81610260565b91508282048414831517610693576106926104b3565b5b5092915050565b610d61806106a75f395ff3fe608060405234801561000f575f5ffd5b5060043610610091575f3560e01c8063313ce56711610064578063313ce5671461013157806370a082311461014f57806395d89b411461017f578063a9059cbb1461019d578063dd62ed3e146101cd57610091565b806306fdde0314610095578063095ea7b3146100b357806318160ddd146100e357806323b872dd14610101575b5f5ffd5b61009d6101fd565b6040516100aa9190610934565b60405180910390f35b6100cd60048036038101906100c891906109e5565b610288565b6040516100da9190610a3d565b60405180910390f35b6100eb610375565b6040516100f89190610a65565b60405180910390f35b61011b60048036038101906101169190610a7e565b61037b565b6040516101289190610a3d565b60405180910390f35b61013961065b565b6040516101469190610ae9565b60405180910390f35b61016960048036038101906101649190610b02565b61066d565b6040516101769190610a65565b60405180910390f35b610187610682565b6040516101949190610934565b60405180910390f35b6101b760048036038101906101b291906109e5565b61070e565b6040516101c49190610a3d565b60405180910390f35b6101e760048036038101906101e29190610b2d565b6108a4565b6040516101f49190610a65565b60405180910390f35b5f805461020990610b98565b80601f016020809104026020016040519081016040528092919081815260200182805461023590610b98565b80156102805780601f1061025757610100808354040283529160200191610280565b820191905f5260205f20905b81548152906001019060200180831161026357829003601f168201915b505050505081565b5f8160055f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925846040516103639190610a65565b60405180910390a36001905092915050565b60035481565b5f8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156103fc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103f390610c12565b60405180910390fd5b8160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205410156104b7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104ae90610c7a565b60405180910390fd5b8160045f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546105039190610cc5565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546105569190610cf8565b925050819055508160055f8673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546105e49190610cc5565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516106489190610a65565b60405180910390a3600190509392505050565b60025f9054906101000a900460ff1681565b6004602052805f5260405f205f915090505481565b6001805461068f90610b98565b80601f01602080910402602001604051908101604052809291908181526020018280546106bb90610b98565b80156107065780601f106106dd57610100808354040283529160200191610706565b820191905f5260205f20905b8154815290600101906020018083116106e957829003601f168201915b505050505081565b5f8160045f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054101561078f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161078690610c12565b60405180910390fd5b8160045f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f8282546107db9190610cc5565b925050819055508160045f8573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461082e9190610cf8565b925050819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516108929190610a65565b60405180910390a36001905092915050565b6005602052815f5260405f20602052805f5260405f205f91509150505481565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f601f19601f8301169050919050565b5f610906826108c4565b61091081856108ce565b93506109208185602086016108de565b610929816108ec565b840191505092915050565b5f6020820190508181035f83015261094c81846108fc565b905092915050565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61098182610958565b9050919050565b61099181610977565b811461099b575f5ffd5b50565b5f813590506109ac81610988565b92915050565b5f819050919050565b6109c4816109b2565b81146109ce575f5ffd5b50565b5f813590506109df816109bb565b92915050565b5f5f604083850312156109fb576109fa610954565b5b5f610a088582860161099e565b9250506020610a19858286016109d1565b9150509250929050565b5f8115159050919050565b610a3781610a23565b82525050565b5f602082019050610a505f830184610a2e565b92915050565b610a5f816109b2565b82525050565b5f602082019050610a785f830184610a56565b92915050565b5f5f5f60608486031215610a9557610a94610954565b5b5f610aa28682870161099e565b9350506020610ab38682870161099e565b9250506040610ac4868287016109d1565b9150509250925092565b5f60ff82169050919050565b610ae381610ace565b82525050565b5f602082019050610afc5f830184610ada565b92915050565b5f60208284031215610b1757610b16610954565b5b5f610b248482850161099e565b91505092915050565b5f5f60408385031215610b4357610b42610954565b5b5f610b508582860161099e565b9250506020610b618582860161099e565b9150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f6002820490506001821680610baf57607f821691505b602082108103610bc257610bc1610b6b565b5b50919050565b7f496e73756666696369656e742062616c616e63650000000000000000000000005f82015250565b5f610bfc6014836108ce565b9150610c0782610bc8565b602082019050919050565b5f6020820190508181035f830152610c2981610bf0565b9050919050565b7f4e6f7420617574686f72697a65640000000000000000000000000000000000005f82015250565b5f610c64600e836108ce565b9150610c6f82610c30565b602082019050919050565b5f6020820190508181035f830152610c9181610c58565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f610ccf826109b2565b9150610cda836109b2565b9250828203905081811115610cf257610cf1610c98565b5b92915050565b5f610d02826109b2565b9150610d0d836109b2565b9250828201905080821115610d2557610d24610c98565b5b9291505056fea2646970667358221220a697a29d6c422ad873c0904ec3dcabe69231988735bf7b41ab5f03b407f7c77464736f6c634300081e0033",
}

// SolTokenABI is the input ABI used to generate the binding from.
// Deprecated: Use SolTokenMetaData.ABI instead.
var SolTokenABI = SolTokenMetaData.ABI

// SolTokenBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SolTokenMetaData.Bin instead.
var SolTokenBin = SolTokenMetaData.Bin

// DeploySolToken deploys a new Ethereum contract, binding an instance of SolToken to it.
func DeploySolToken(auth *bind.TransactOpts, backend bind.ContractBackend, initialSupply *big.Int) (common.Address, *types.Transaction, *SolToken, error) {
	parsed, err := SolTokenMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SolTokenBin), backend, initialSupply)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SolToken{SolTokenCaller: SolTokenCaller{contract: contract}, SolTokenTransactor: SolTokenTransactor{contract: contract}, SolTokenFilterer: SolTokenFilterer{contract: contract}}, nil
}

// SolToken is an auto generated Go binding around an Ethereum contract.
type SolToken struct {
	SolTokenCaller     // Read-only binding to the contract
	SolTokenTransactor // Write-only binding to the contract
	SolTokenFilterer   // Log filterer for contract events
}

// SolTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type SolTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SolTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SolTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SolTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SolTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SolTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SolTokenSession struct {
	Contract     *SolToken         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SolTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SolTokenCallerSession struct {
	Contract *SolTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// SolTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SolTokenTransactorSession struct {
	Contract     *SolTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SolTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type SolTokenRaw struct {
	Contract *SolToken // Generic contract binding to access the raw methods on
}

// SolTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SolTokenCallerRaw struct {
	Contract *SolTokenCaller // Generic read-only contract binding to access the raw methods on
}

// SolTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SolTokenTransactorRaw struct {
	Contract *SolTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSolToken creates a new instance of SolToken, bound to a specific deployed contract.
func NewSolToken(address common.Address, backend bind.ContractBackend) (*SolToken, error) {
	contract, err := bindSolToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SolToken{SolTokenCaller: SolTokenCaller{contract: contract}, SolTokenTransactor: SolTokenTransactor{contract: contract}, SolTokenFilterer: SolTokenFilterer{contract: contract}}, nil
}

// NewSolTokenCaller creates a new read-only instance of SolToken, bound to a specific deployed contract.
func NewSolTokenCaller(address common.Address, caller bind.ContractCaller) (*SolTokenCaller, error) {
	contract, err := bindSolToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SolTokenCaller{contract: contract}, nil
}

// NewSolTokenTransactor creates a new write-only instance of SolToken, bound to a specific deployed contract.
func NewSolTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*SolTokenTransactor, error) {
	contract, err := bindSolToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SolTokenTransactor{contract: contract}, nil
}

// NewSolTokenFilterer creates a new log filterer instance of SolToken, bound to a specific deployed contract.
func NewSolTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*SolTokenFilterer, error) {
	contract, err := bindSolToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SolTokenFilterer{contract: contract}, nil
}

// bindSolToken binds a generic wrapper to an already deployed contract.
func bindSolToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SolTokenMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SolToken *SolTokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SolToken.Contract.SolTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SolToken *SolTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SolToken.Contract.SolTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SolToken *SolTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SolToken.Contract.SolTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SolToken *SolTokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SolToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SolToken *SolTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SolToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SolToken *SolTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SolToken.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_SolToken *SolTokenCaller) Allowance(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SolToken.contract.Call(opts, &out, "allowance", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_SolToken *SolTokenSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _SolToken.Contract.Allowance(&_SolToken.CallOpts, arg0, arg1)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_SolToken *SolTokenCallerSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _SolToken.Contract.Allowance(&_SolToken.CallOpts, arg0, arg1)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_SolToken *SolTokenCaller) BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SolToken.contract.Call(opts, &out, "balanceOf", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_SolToken *SolTokenSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _SolToken.Contract.BalanceOf(&_SolToken.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_SolToken *SolTokenCallerSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _SolToken.Contract.BalanceOf(&_SolToken.CallOpts, arg0)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_SolToken *SolTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _SolToken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_SolToken *SolTokenSession) Decimals() (uint8, error) {
	return _SolToken.Contract.Decimals(&_SolToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_SolToken *SolTokenCallerSession) Decimals() (uint8, error) {
	return _SolToken.Contract.Decimals(&_SolToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_SolToken *SolTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _SolToken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_SolToken *SolTokenSession) Name() (string, error) {
	return _SolToken.Contract.Name(&_SolToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_SolToken *SolTokenCallerSession) Name() (string, error) {
	return _SolToken.Contract.Name(&_SolToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_SolToken *SolTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _SolToken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_SolToken *SolTokenSession) Symbol() (string, error) {
	return _SolToken.Contract.Symbol(&_SolToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_SolToken *SolTokenCallerSession) Symbol() (string, error) {
	return _SolToken.Contract.Symbol(&_SolToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_SolToken *SolTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SolToken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_SolToken *SolTokenSession) TotalSupply() (*big.Int, error) {
	return _SolToken.Contract.TotalSupply(&_SolToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_SolToken *SolTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _SolToken.Contract.TotalSupply(&_SolToken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_SolToken *SolTokenTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _SolToken.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_SolToken *SolTokenSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _SolToken.Contract.Approve(&_SolToken.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_SolToken *SolTokenTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _SolToken.Contract.Approve(&_SolToken.TransactOpts, spender, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_SolToken *SolTokenTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _SolToken.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_SolToken *SolTokenSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _SolToken.Contract.Transfer(&_SolToken.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_SolToken *SolTokenTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _SolToken.Contract.Transfer(&_SolToken.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_SolToken *SolTokenTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _SolToken.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_SolToken *SolTokenSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _SolToken.Contract.TransferFrom(&_SolToken.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_SolToken *SolTokenTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _SolToken.Contract.TransferFrom(&_SolToken.TransactOpts, from, to, value)
}

// SolTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the SolToken contract.
type SolTokenApprovalIterator struct {
	Event *SolTokenApproval // Event containing the contract specifics and raw log

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
func (it *SolTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SolTokenApproval)
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
		it.Event = new(SolTokenApproval)
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
func (it *SolTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SolTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SolTokenApproval represents a Approval event raised by the SolToken contract.
type SolTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_SolToken *SolTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*SolTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _SolToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &SolTokenApprovalIterator{contract: _SolToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_SolToken *SolTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *SolTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _SolToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SolTokenApproval)
				if err := _SolToken.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_SolToken *SolTokenFilterer) ParseApproval(log types.Log) (*SolTokenApproval, error) {
	event := new(SolTokenApproval)
	if err := _SolToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SolTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the SolToken contract.
type SolTokenTransferIterator struct {
	Event *SolTokenTransfer // Event containing the contract specifics and raw log

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
func (it *SolTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SolTokenTransfer)
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
		it.Event = new(SolTokenTransfer)
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
func (it *SolTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SolTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SolTokenTransfer represents a Transfer event raised by the SolToken contract.
type SolTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_SolToken *SolTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SolTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SolToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &SolTokenTransferIterator{contract: _SolToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_SolToken *SolTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *SolTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SolToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SolTokenTransfer)
				if err := _SolToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_SolToken *SolTokenFilterer) ParseTransfer(log types.Log) (*SolTokenTransfer, error) {
	event := new(SolTokenTransfer)
	if err := _SolToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
