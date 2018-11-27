package pricebook

import (
	"financialExchange/sql"
	"financialExchange/util"
	"time"
	"financialExchange/model"
	"math"
)

type TimePeriod int

const(
	Day 	TimePeriod = 0
	Week 	TimePeriod = 1
	Month 	TimePeriod = 2
	ThreeMonth TimePeriod = 3
	Year TimePeriod = 4
	All TimePeriod = 5
)


type PriceController struct{
	Database *sql.MySqlDB
	Logger *util.Logger
}

func NewPriceController(database *sql.MySqlDB, logger *util.Logger)(*PriceController){
	return &PriceController{
		Database: database,
		Logger: logger,
	}
}

func (self *PriceController)WeightedAverageSharePrice(start int64, end int64, security int64)(float64, error){
	matchingTransactions, err := self.Database.GetAllTransactionsForTimePeriodForSecurity(start, end, security)
	if err != nil{
		return 0, err
	}
	if len(matchingTransactions) == 0{
		return -1, nil
	}

	numSharesTotal := 0.0
	summedSharePrice := 0.0

	for _, transaction := range(matchingTransactions){
		numSharesTotal += float64(transaction.NumShares)

		transactionCostPerShare, _ := transaction.CostPerShare.Float64()
		costPerShare := math.Floor(transactionCostPerShare*100)/100

		//if !exact{
		//	return 0, errors.New("UnexactTransactionCostPerShare")
		//}

		summedSharePrice += float64(transaction.NumShares) * costPerShare
	}

	return summedSharePrice/numSharesTotal, nil
}

func (self *PriceController)GetSecurityChartForTimePeriod(timePeriod TimePeriod, security int64)([]model.PricePoint, error){
	var startTime time.Time
	var timeIncrement int64
	var numTicks int
	endTime := time.Now()

	switch timePeriod{
	//5 mim intervals
	case Day:
		startTime = endTime.Add(-1 * 24* time.Hour)
		numTicks = 288
	//5 min intervals
	case Week:
		startTime = endTime.Add(-1 * 24 * 7 * time.Hour)
		numTicks = 2016
	//Day intervals
	case Month:
		startTime = endTime.Add(-1 * 24 * 31 * time.Hour)
		numTicks = 31
	//Day intervals
	case ThreeMonth:
		startTime = endTime.Add(-1 * 24 * 92 * time.Hour)
		numTicks = 92
	//Day intervals
	case Year:
		startTime = endTime.Add(-1 * 24 * 365 * time.Hour)
		numTicks = 365
	case All:
		security, err := self.Database.GetSecurityByID(security)
		if err != nil{
			return []model.PricePoint{}, err
		}
		startTime = time.Unix(security.Created,0)
		numTicks = int((endTime.Unix()-startTime.Unix()) / (60 * 60 * 24))
	}

	timeIncrement = int64(endTime.Sub(startTime).Seconds()/float64(numTicks))

	prices := make([]model.PricePoint, 0)

	for tick := int64(0); tick < int64(numTicks); tick += 1{
		startPeriod := int64(tick*timeIncrement)+startTime.Unix()
		endPeriod := int64((tick+1)*timeIncrement)+startTime.Unix()
		WASP, err := self.WeightedAverageSharePrice(startPeriod, endPeriod, security)
		if err != nil{
			return []model.PricePoint{}, err
		}else{
			if WASP == float64(-1){
				if len(prices) != 0{
					pricePoint := model.PricePoint{
						TimeStamp: endPeriod,
						PricePoint: prices[len(prices)-1].PricePoint,
					}
					prices = append(prices, pricePoint)
				}else{
					tickBack := int64(1)
					securityModel, err := self.Database.GetSecurityByID(security)
					if err != nil{
						return []model.PricePoint{}, err
					}
					entity, err := self.Database.GetEntityByID(securityModel.Entity)
					if err != nil{
						return []model.PricePoint{}, err
					}
					IPOOrder, err := self.Database.GetTransactionByID(entity.IPO)
					if err != nil{
						return []model.PricePoint{}, err
					}

					for WASP == -1 && startPeriod-(tickBack*timeIncrement) > IPOOrder.Created{
						WASP, err := self.WeightedAverageSharePrice(startPeriod-(tickBack*timeIncrement),endPeriod-(tickBack*timeIncrement), security)
						if err != nil{
							return []model.PricePoint{}, err
						}
						tickBack += 1
						if WASP != -1{
							pricePoint := model.PricePoint{
								TimeStamp: endPeriod,
								PricePoint: WASP,
							}
							prices = append(prices, pricePoint)
						}
					}


				}
			}else{
				pricePoint := model.PricePoint{
					TimeStamp: endPeriod,
					PricePoint: WASP,
				}
				prices = append(prices, pricePoint)
			}

		}
	}
	return prices, nil
}

func(self *PriceController)GetCurrPriceOfSecurity(securityID int64)(model.PricePoint, error){
	security, err := self.Database.GetSecurityByID(securityID)
	if err != nil{
		return model.PricePoint{}, err
	}
	WASP := float64(-1)
	tickBack := int64(0)
	entity, err := self.Database.GetEntityByID(security.Entity)
	if err != nil{
		return model.PricePoint{}, err
	}
	IPOOrder, err := self.Database.GetTransactionByID(entity.IPO)
	if err != nil{
		return model.PricePoint{}, err
	}
	startPeriod := time.Now().Add(time.Minute * -5).Unix()
	endPeriod := time.Now().Unix()

	for WASP == -1 && endPeriod-(tickBack*300) > IPOOrder.Created{
		WASP, err := self.WeightedAverageSharePrice(startPeriod-(tickBack*300),endPeriod-(tickBack*300), security.Id)
		if err != nil{
			return model.PricePoint{}, err
		}
		tickBack += 1
		if WASP != -1{
			return model.PricePoint{
				TimeStamp: time.Now().Unix(),
				PricePoint: WASP,
			}, nil
		}
	}
	return model.PricePoint{}, nil
}