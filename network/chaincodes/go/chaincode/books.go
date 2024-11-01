package chaincode

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
)

// define struct of smart contract
type SmartContract struct {
	contractapi.Contract
}

// registration
func (s *SmartContract) UserRegistration(
	_ctx contractapi.TransactionContextInterface,
	_userID string,
	_userType string,
	_realInfoHash string) error {
	user := User{
		UserID:       _userID,
		UserType:     _userType,
		RealInfoHash: _realInfoHash,
		BookList:     []*Book{},
	}
	userAsBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = _ctx.GetStub().PutState(_userID, userAsBytes)
	if err != nil {
		return err
	}
	return nil
}

// upload books onto chain
// this function passes in basic info of books
func (s *SmartContract) UploadOntoChain(
	_ctx contractapi.TransactionContextInterface,
	_userID string,
	_traceability_code string,
	_arg1 string,
	_arg2 string,
	_arg3 string,
	_arg4 string,
	_arg5 string,
	_arg6 string,
	_arg7 string) (string, error) {
	// get user type, if fail return error message
	userType, err := s.GetUserType(_ctx, _userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user type: %v", err)
	}

	// get book data according to the traceability_code
	bookAsBytes, err := _ctx.GetStub().GetState(_traceability_code)
	if err != nil {
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}

	// transfer books info into book struct
	var book Book
	if bookAsBytes != nil {
		err = json.Unmarshal(bookAsBytes, &book)
		if err != nil {
			return "", fmt.Errorf("failed to Unmarshal book: %v", err)
		}
	}

	// get timestamp
	txtime, err := _ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return "", fmt.Errorf("failed to read TxTimestamp: %v", err)
	}
	timeLocation, _ := time.LoadLocation("Asia/Shanghai")
	time := time.Unix(txtime.Seconds, 0).In(timeLocation).Format("2006-01-02 15:04:05")

	// acquire txid
	txid := _ctx.GetStub().GetTxID()

	// apply traceable codes to books
	book.Traceability_code = _traceability_code

	// according to different userType, upload different info onto chain
	switch userType {
	case "author":
		// this struct is the info struct of books as a work of art, for authors
		book.Author_input.Author_bookTitle = _arg1
		book.Author_input.Author_bookAuthors = _arg2
		book.Author_input.Author_bookPublishTime = _arg3
		book.Author_input.Author_bookPublisher = _arg4
		book.Author_input.Author_bookEditionNo = _arg5
		book.Author_input.Author_bookPrintedNum = _arg6
		book.Author_input.Author_bookISBN = _arg7
		book.Author_input.Author_TxID = txid
		book.Author_input.Author_Timestamp = time

	case "press":
		// this struct is the info struct of books as a product for printing, for press
		book.Press_input.Press_bookTitle = _arg1
		book.Press_input.Press_bookAuthors = _arg2
		book.Press_input.Press_bookPublishTime = _arg3
		book.Press_input.Press_bookEditionNo = _arg4
		book.Press_input.Press_bookISBN = _arg5
		book.Press_input.Press_bookPrintTime = _arg6
		book.Press_input.Press_bookPrintedNum = _arg7
		book.Press_input.Press_TxID = txid
		book.Press_input.Press_Timestamp = time

	case "Logistics":
		// this struct is the info struct of books as an item for delivering, for logistics
		book.Logistics_input.Logistics_bookTitle = _arg1
		book.Logistics_input.Logistics_bookAuthors = _arg2
		book.Logistics_input.Logistics_bookISBN = _arg3
		book.Logistics_input.Logistics_bookDeparture = _arg4
		book.Logistics_input.Logistics_bookDestination = _arg5
		book.Logistics_input.Logistics_bookStorageShelf = _arg6
		book.Logistics_input.Logistics_bookDelivererID = _arg7
		book.Logistics_input.Logistics_TxID = txid
		book.Logistics_input.Logistics_Timestamp = time

	case "Retailer":
		// this struct is the info struct of books as a good for selling, for retailers
		book.Retailer_input.Retailer_bookISBN = _arg1
		book.Retailer_input.Retailer_bookCostPerUnit = _arg2
		book.Retailer_input.Retailer_bookPrice = _arg3
		book.Retailer_input.Retailer_bookPurchaseFrom = _arg4
		book.Retailer_input.Retailer_bookArriveTime = _arg5
		book.Retailer_input.Retailer_bookOnShelfNum = _arg6
		book.Retailer_input.Retailer_bookRetailerContact = _arg7
		book.Retailer_input.Retailer_TxID = txid
		book.Retailer_input.Retailer_Timestamp = time
	}

	// transfer struct type into json
	bookAsBytes, err = json.Marshal(book)
	if err != nil {
		return "", fmt.Errorf("failed to marshal fruit: %v", err)
	}

	// upload the data onto chain
	err = _ctx.GetStub().PutState(_traceability_code, bookAsBytes)
	if err != nil {
		return "", fmt.Errorf("failed to push fruit: %v", err)
	}

	// save book info into book list
	err = s.AddBook(_ctx, _userID, &book)
	if err != nil {
		return "", fmt.Errorf("failed to add book to user: %v", err)
	}

	return txid, nil
}

// add books to user's book lists
func (s *SmartContract) AddBook(
	_ctx contractapi.TransactionContextInterface,
	_userID string,
	_book *Book) error {
	// read user list from chain
	userBytes, err := _ctx.GetStub().GetState(_userID)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if userBytes == nil {
		return fmt.Errorf("user %s does not exist", _userID)
	}

	// transfer result into user struct
	var user User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return err
	}

	user.BookList = append(user.BookList, _book)
	userAsBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = _ctx.GetStub().PutState(_userID, userAsBytes)
	if err != nil {
		return err
	}
	return nil
}

// acquire user type
func (s *SmartContract) GetUserType(
	_ctx contractapi.TransactionContextInterface,
	_userID string) (string, error) {
	userBytes, err := _ctx.GetStub().GetState(_userID)
	if err != nil {
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	if userBytes == nil {
		return "", fmt.Errorf("user %s does not exist", _userID)
	}

	// transfer result into user struct
	var user User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return "", err
	}
	return user.UserType, nil
}

// acquire user info
func (s *SmartContract) GetUserInfo(_ctx contractapi.TransactionContextInterface, _userID string) (*User, error) {
	userBytes, err := _ctx.GetStub().GetState(_userID)
	if err != nil {
		return &User{}, fmt.Errorf("failed to read from world state: %v", err)
	}
	if userBytes == nil {
		return &User{}, fmt.Errorf("the user %s does not exist", _userID)
	}
	// transfer result into user struct
	var user User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return &User{}, err
	}
	return &user, nil
}

// acquire on chain books data
func (s *SmartContract) GetBookInfo(_ctx contractapi.TransactionContextInterface, _traceability_code string) (*Book, error) {
	BookAsBytes, err := _ctx.GetStub().GetState(_traceability_code)
	if err != nil {
		return &Book{}, fmt.Errorf("failed to read from world state: %v", err)
	}

	// transfer result into book struct
	var book Book
	if BookAsBytes != nil {
		err = json.Unmarshal(BookAsBytes, &book)
		if err != nil {
			return &Book{}, fmt.Errorf("failed to Unmarshal book: %v", err)
		}
		if book.Traceability_code != "" {
			return &book, nil
		}
	}
	return &Book{}, fmt.Errorf("the book %s does not exist", _traceability_code)
}

// acquire book id list
func (s *SmartContract) GetBookList(_ctx contractapi.TransactionContextInterface, _userID string) ([]*Book, error) {
	userBytes, err := _ctx.GetStub().GetState(_userID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if userBytes == nil {
		return nil, fmt.Errorf("the user %s does not exist", _userID)
	}
	// transfer result into user struct
	var user User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return nil, err
	}
	return user.BookList, nil
}

// acquire all book info
func (s *SmartContract) GetAllBookInfo(_ctx contractapi.TransactionContextInterface) ([]Book, error) {
	bookListAsBytes, err := _ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	defer bookListAsBytes.Close()
	var books []Book
	for bookListAsBytes.HasNext() {
		queryResponse, err := bookListAsBytes.Next()
		if err != nil {
			return nil, err
		}
		var book Book
		err = json.Unmarshal(queryResponse.Value, &book)
		if err != nil {
			return nil, err
		}
		// filt non-book data
		if book.Traceability_code != "" {
			books = append(books, book)
		}
	}
	return books, nil
}

// acquire upload book history
func (s *SmartContract) GetBookHistory(_ctx contractapi.TransactionContextInterface, _traceability_code string) ([]HistoryQueryResult, error) {
	log.Printf("GetAssetHistory: ID %v", _traceability_code)

	resultsIterator, err := _ctx.GetStub().GetHistoryForKey(_traceability_code)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var book Book
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &book)
			if err != nil {
				return nil, err
			}
		} else {
			book = Book{
				Traceability_code: _traceability_code,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}
		// specify target time zone
		targetLocation, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			return nil, err
		}

		// transfer timestamp into target time zone
		timestamp = timestamp.In(targetLocation)
		// standardise timestamp into required format
		formattedTime := timestamp.Format("2006-01-02 15:04:05")

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: formattedTime,
			Record:    &book,
			IsDeleted: response.IsDelete,
		}
		records = append(records, record)
	}

	return records, nil
}
