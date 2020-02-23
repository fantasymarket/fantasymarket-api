package main

type Stocks struct {
	index_value int
	category string
	
}

// TODO: DECLARE INTERFACE WITH COMMON METHODS FOR THE STOCKS TO DETERMINE IN WHICH CATEGORY THEY FALL IN, BECAUSE
// TODO: DEPENDING ON WHICH CATEGORY THEY'RE IN THE NUMBERS WILL REACT DIFFERENTLY TO THE EVENTS

// Information Technology

var Aapl = Stocks{index_value: 0, category: "Tech"}
var Msft = Stocks{index_value: 0, category: "Tech"}
var Fb = Stocks{index_value: 0, category: "Tech"}
var Nflx = Stocks{index_value: 0, category: "Tech"}
var Amzn = Stocks{index_value: 0, category: "Tech"}
var Tcehy = Stocks{index_value: 0, category: "Tech"}
var Ahlh = Stocks{index_value: 0, category: "Tech"}
var Goog = Stocks{index_value: 0, category: "Tech"}


// Finance

var Axp = Stocks{index_value: 0, category: "Finance"}
var Jpm = Stocks{index_value: 0, category: "Finance"}
var Gs = Stocks{index_value: 0, category: "Finance"}
var Db = Stocks{index_value: 0, category: "Finance"}


// Energy

var Rds_a = Stocks{index_value: 0, category: "Energy"}
var Tsla = Stocks{index_value: 0, category: "Energy"}
var Nep = Stocks{index_value: 0, category: "Energy"}
var Fslr = Stocks{index_value: 0, category: "Energy"}


// Healthcare

var Idxx = Stocks{index_value: 0, category: "Healthcare"}
var Antm = Stocks{index_value: 0, category: "Healthcare"}
var Dva = Stocks{index_value: 0, category: "Healthcare"}


// Raw Materials

var Gld = Stocks{index_value: 0, category: "Raw Materials"}

var Stock_map = map[string]*Stocks{
	"Aapl": &Aapl,
	"Msft": &Msft,
	"Fb": &Fb,
	"Nflx": &Nflx,
	"Amzn": &Amzn,
	"Tcehy": &Tcehy,
	"Ahlh": &Ahlh,
	"Goog": &Goog,
	"Axp": &Axp,
	"Jpm": &Jpm,
	"Gs": &Gs,
	"Db": &Db,
	"Rds_a": &Rds_a,
	"Tsla": &Tsla,
	"Nep": &Nep,
	"Fslr": &Fslr,
	"Idxx": &Idxx,
	"Antm": &Antm,
	"Dva": &Dva,
	"Gld": &Gld,

}


