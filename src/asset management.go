package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SimpleChaincode struct {
}

type User struct {
	User   string `json:"user"`
	Amount string `json:"amount"`
}

const total = "Total"

var userNumber int

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	userNumber = 0
	t.mint(stub, []string{total, "0"})
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "mint" {
		return t.mint(stub, args)
	} else if function == "balanceOf" {
		return t.balanceOf(stub, args)
	} else if function == "withdraw" {
		return t.withdraw(stub, args)
	} else if function == "transfer" {
		return t.transfer(stub, args)
	} else if function == "totalAmount" {
		return t.totalAmount(stub)
	} else if function == "queryAllUsers" {
		return t.queryAllUsers(stub)
	} else if function == "deleteUser" {
		return t.deleteUser(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"mint\" \"transfer\" \"totalAmount\" \"balanceOf\" \"withdraw\" \"queryAllUsers\" \"deleteUser\"")
}

func (t *SimpleChaincode) mint(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	user := args[0]
	amount, _ := strconv.Atoi(args[1])

	if user == "Total" {
		return shim.Error("Invalid user name.")
	}
	if amount < 0 {
		return shim.Error("Amount can't be negative.")
	}

	if err := stub.PutState(user, []byte(strconv.Itoa(amount))); err != nil {
		return shim.Error(err.Error())
	}

	t.addTotalAmount(stub, amount)

	return shim.Success(nil)
}

func (t *SimpleChaincode) addTotalAmount(stub shim.ChaincodeStubInterface, amount int) pb.Response {

	totalAmountAsBytes := t.balanceOf(stub, []string{total})
	totalAmount, _ := strconv.Atoi(string(totalAmountAsBytes))
	totalAmount += amount

	if err := stub.PutState(total, []byte(strconv.Itoa(totalAmount))); err != nil {
		return shim.Error(err.Error())
	}

	return nil
}

func (t *SimpleChaincode) balanceOf(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	user := args[0]

	balanceAsBytes, err := stub.GetState(user)
	if err != nil {
		return shim.Error(err.Error())
	}

	return balanceAsBytes
}

func (t *SimpleChaincode) withdraw(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	user := args[0]
	amount, _ := strconv.Atoi(args[1])

	if amount < 0 {
		return shim.Error("Amount can't be negative.")
	}

	balanceAsBytes := t.balanceOf(stub, []string{user})

	balance, _ := strconv.Atoi(string(balanceAsBytes))
	balance = balance - amount

	if balance < 0 {
		return shim.Error("The amount is more than the balance.")
	}

	if err := stub.PutState(user, []byte(strconv.Itoa(balance))); err != nil {
		return shim.Error(err.Error())
	}

	return nil
}

func (t *SimpleChaincode) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	info := []string{args[0], args[2]}
	t.withdraw(stub, info)

	beneficiary := args[1]

	balance, _ := strconv.Atoi(string(t.balanceOf(stub, []string{beneficiary})))
	amount, _ := strconv.Atoi(args[2])
	balance = balance + amount

	if err := stub.PutState(beneficiary, []byte(strconv.Itoa(balance))); err != nil {
		return shim.Error(err.Error())
	}

	return nil
}

func (t *SimpleChaincode) totalAmount(stub shim.ChaincodeStubInterface) pb.Response {
	totalAmount, _ := strconv.Atoi(string(t.balanceOf(stub, []string{total})))

	fmt.Println("Total amount: %d\n", totalAmount)

	return nil
}

func (t *SimpleChaincode) queryAllUsers(stub shim.ChaincodeStubInterface) pb.Response {
	return nil
}

// Deletes an entity from state
func (t *SimpleChaincode) deleteUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
