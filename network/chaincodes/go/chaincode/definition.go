package chaincode

type User struct {
	UserID       string  `json:"userID"`
	UserType     string  `json:"userType"`
	RealInfoHash string  `json:"realInfoHash"`
	BookList     []*Book `json:"bookList"`
}

type Book struct {
	Traceability_code string          `json:"traceability_code"`
	Author_input      Author_input    `json:"author_input"`
	Press_input       Press_input     `json:"press_input"`
	Logistics_input   Logistics_input `json:"logistics_input"`
	Retailer_input    Retailer_input  `json:"retailer_input"`
}

// HistoryQueryResult structure used for handling result of history query
type HistoryQueryResult struct {
	Record    *Book  `json:"record"`
	TxId      string `json:"txId"`
	Timestamp string `json:"timestamp"`
	IsDeleted bool   `json:"isDeleted"`
}

type Author_input struct {
	Author_bookTitle       string `json:"author_bookTitle"`
	Author_bookAuthors     string `json:"author_bookAuthors"`
	Author_bookPublishTime string `json:"author_bookPublishTime"`
	Author_bookPublisher   string `json:"author_bookPublisher"`
	Author_bookEditionNo   string `json:"author_bookEditionNo"`
	Author_bookPrintedNum  string `json:"author_bookPrintedNum"`
	Author_bookISBN        string `json:"author_bookISBN"`
	Author_TxID            string `json:"author_txID"`
	Author_Timestamp       string `json:"author_timestamp"`
}

type Press_input struct {
	Press_bookTitle       string `json:"press_bookTitle"`
	Press_bookAuthors     string `json:"press_bookAuthors"`
	Press_bookPublishTime string `json:"press_bookPublishTime"`
	Press_bookEditionNo   string `json:"press_bookEditionNo"`
	Press_bookISBN        string `json:"press_bookISBN"`
	Press_bookPrintTime   string `json:"press_bookPrintTime"`
	Press_bookPrintedNum  string `json:"press_bookPrintedNum"`
	Press_TxID            string `json:"press_txID"`
	Press_Timestamp       string `json:"press_timestamp"`
}

type Logistics_input struct {
	Logistics_bookTitle        string `json:"logistics_bookTitle"`
	Logistics_bookAuthors      string `json:"logistics_bookAuthors"`
	Logistics_bookISBN         string `json:"logistics_bookISBN"`
	Logistics_bookDeparture    string `json:"logistics_bookDeparture"`
	Logistics_bookDestination  string `json:"logistics_bookDestination"`
	Logistics_bookStorageShelf string `json:"logistics_bookStorageShelf"`
	Logistics_bookDelivererID  string `json:"logistics_bookDelivererID"`
	Logistics_TxID             string `json:"logistics_txID"`
	Logistics_Timestamp        string `json:"logistics_timestamp"`
}

type Retailer_input struct {
	Retailer_bookISBN            string `json:"retailer_bookISBN"`
	Retailer_bookCostPerUnit     string `json:"retailer_bookCostPerUnit"`
	Retailer_bookPrice           string `json:"retailer_bookPrice"`
	Retailer_bookPurchaseFrom    string `json:"retailer_bookPurchaseFrom"`
	Retailer_bookArriveTime      string `json:"retailer_bookArriveTime"`
	Retailer_bookOnShelfNum      string `json:"retailer_bookOnShelfNum"`
	Retailer_bookRetailerContact string `json:"retailer_bookRetailerContact"`
	Retailer_TxID                string `json:"retailer_TxID"`
	Retailer_Timestamp           string `json:"retailer_Timestamp"`
}
