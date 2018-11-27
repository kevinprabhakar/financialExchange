var currAction = 0
var currUser = {}
var currSecurity = 1
var currTimePeriod = 0
var currentPrice = 0.0
var currChart = null
var currSymbol = null


$(document).ready(function(){
      var BearerString = 'Bearer ' + api.readCookie("accessToken")
      GetSymbolOfCurrSecurity()
      GenerateNewPriceLine()
      $.ajax({
        url: "/api/customer",
        headers: {
            'Authorization': BearerString
        },
        type: "GET",
        contentType: "application/json",
        success: function(response, textStatus, jQxhr){
               parsedJSON = JSON.parse(response)
               currUser = parsedJSON
        },
        error: function( jqXhr, textStatus, errorThrown ){
                console.log( errorThrown );
        }
      })
      $("#securitySearch").keyup(function(){
            securityName = document.getElementById("securitySearch").value;
            if (securityName.length === 0){
                return
            }

            $.ajax({
                url: "/api/search",
                headers: {
                    'Authorization': BearerString
                },
                contentType: 'application/json',
                type: "POST",
                data: JSON.stringify({"prefix":securityName}),
                success: function(response, textStatus, jQxhr){
                       searchResults = JSON.parse(response)
                       newHTMLString = ""
                       for (var i = 0; i < searchResults.length; i++){
                            newHTMLString += '<option value="'+searchResults[i].entityName+" ("+searchResults[i].symbol+')">'
                       }
                       var secList = document.getElementById("securitiesList")
                       secList.innerHTML = newHTMLString
                },
                error: function( jqXhr, textStatus, errorThrown ){
                        console.log( errorThrown );
                }

            });

        });
      $.ajax({
              url: "/api/currprice/"+currSecurity.toString(),
              type: "GET",
              contentType: "application/json",
              success: function(response, textStatus, jQxhr){
                     parsedJSON = JSON.parse(response)
                     currentPrice = parsedJSON.pricePoint
                     var currPrice = document.getElementById("currPrice")
                     currPrice.textContent = "Current Price: $" + currentPrice.toString()

              },
              error: function( jqXhr, textStatus, errorThrown ){
                      console.log( errorThrown );
              }
            })
      $.ajax({
              url: "/api/customer/portfolio",
              headers: {
                  'Authorization': BearerString
              },
              type: "GET",
              contentType: "application/json",
              success: function(response, textStatus, jQxhr){
                     currUserPortfolio = JSON.parse(response)
                     newHTMLString = "<table><tr><td>Cash Value: $</td><td>" + currUserPortfolio["cashValue"] + "</td></tr><tr><td><a href='ownedShares.html'>Stock Value</a>: $</td><td>"+currUserPortfolio["stockValue"]+"</td></tr></table>"


                     var userPortfolio = document.getElementById("currUserPortfolio")
                     userPortfolio.innerHTML = newHTMLString
              },
              error: function( jqXhr, textStatus, errorThrown ){
                      console.log( errorThrown );
              }
            })

      $.ajax({
        url: "/api/customer/orders",
        headers: {
            'Authorization': BearerString
        },
        type: "GET",
        contentType: "application/json",
        success: function(response, textStatus, jQxhr){
           orderList = JSON.parse(response);
           newInnerHTML = "";
           console.log(orderList);
           if (orderList.length != 0){
                for (var i = orderList.length-1; i >= 0; i--){
                    newHTMLString = ""
                    newHTMLString += "Symbol: " + orderList[i]["symbol"] + "<br>"
                    newHTMLString += "Initial Num Shares: " + orderList[i]["numShares"] + "<br>"
                    newHTMLString += "Num Shares Remaining: " + orderList[i]["numSharesRemaining"] + "<br>"
                    newHTMLString += "Time Placed: " + new moment.unix(orderList[i]['created']).format('MM/DD/YYYY HH:MM') + "<br>"
                    if (orderList[i]["orderStatus"] === 2){
                        newHTMLString += "Time Fulfilled: " + new moment.unix(orderList[i]['fulfilled']).format('MM/DD/YYYY HH:MM') + "<br>"
                    }
                    if (orderList[i]["investorAction"] === 0){
                        newHTMLString += "Investor Action: Buy <br>"
                        newHTMLString += "Max Cost Per Share: $" + orderList[i]["costPerShare"] + "<br>"

                    }else{
                        newHTMLString += "Investor Action: Sell <br>"
                        newHTMLString += "Min Cost Per Share: $" + orderList[i]["costPerShare"] + "<br>"

                    }
                    newHTMLString += "<br>"

                    newInnerHTML += newHTMLString
                }
                var currOrders = document.getElementById("currUserOrders")
                currOrders.innerHTML = newInnerHTML
           }
        },
        error: function( jqXhr, textStatus, errorThrown ){
                console.log( errorThrown );
        }
      })


});



function GenerateNewPriceLine(){
    var BearerString = 'Bearer ' + api.readCookie("accessToken")

    $.ajax({
          url: "/api/prices",
          type: "POST",
          contentType: 'application/json',
          data: JSON.stringify({"timePeriod":currTimePeriod, "security":currSecurity}),
          success: function(response, textStatus, jQxhr){
              chartData = JSON.parse(response);
              pricePoints = []
              timeStamps = []
              for (var i = 0; i < chartData.length; i++){
                  pricePoints.push(chartData[i]['pricePoint'])
                  timeStamps.push(new moment.unix(chartData[i]['timeStamp']).format('MM/DD/YYYY HH:MM:SS'))
              }

              var chartData = {
                labels: timeStamps,
                datasets: [{
                  label: currSymbol,
                  data: pricePoints,
                }]
              };

              // Create a new line chart object where as first parameter we pass in a selector
              // that is resolving to our chart container element. The Second parameter
              // is the actual data object.
              var chLine = document.getElementById('priceLine');

              if (chLine) {
                currChart = new Chart(chLine, {
                type: 'line',
                data: chartData,
                showArea: true,
                lineSmooth: false
                });
              }
          },
          error: function( jqXhr, textStatus, errorThrown ){
                  console.log( errorThrown );
          }
      });
}

function changePriceLineWindow(evt, window, elemID){
    tablinks = document.getElementsByClassName("tablinks2");
    for (i = 0; i < tablinks.length; i++) {
        tablinks[i].className = tablinks[i].className.replace(" active", "");
    }
    var currButton = document.getElementById(elemID)
    evt.currentTarget.className += " active";
    currTimePeriod = window
    GenerateNewPriceLine()
}

function updateTotalAmount(evt) {
    numShares = parseInt(document.getElementById("numShares").value);
    amountPerShare = parseFloat(document.getElementById("amountPerShare").value);
    if (!isNaN(numShares)&&!isNaN(amountPerShare)){
        console.log(numShares*amountPerShare);
        var totalAmount = document.getElementById("totalAmount");
        console.log(toString(numShares*amountPerShare));
        totalAmount.textContent = "  $"+ numShares*amountPerShare;
    }
}

function changeAction(evt, action) {
    tablinks = document.getElementsByClassName("tablinks");
    for (i = 0; i < tablinks.length; i++) {
        tablinks[i].className = tablinks[i].className.replace(" active", "");
    }
    var currButton = document.getElementById(action)
    evt.currentTarget.className += " active";

    if (currButton.id==="Sell"){
        currAction = 1;
    }else{
        currAction = 0;
    }

    console.log(currAction);
}

function getAppropriateSecurities(event){
    searchParam = document.getElementById("securitySearch").value;
    var BearerString = 'Bearer ' + api.readCookie("accessToken")

    init = searchParam.indexOf('(');
    fin = searchParam.indexOf(')');
    realSecName = searchParam.substr(init+1,fin-init-1);

    $.ajax({
        url: '/api/security/'+realSecName,
        headers: {
            'Authorization': BearerString
        },
        contentType: 'application/json',
        type: "GET",
        success: function(response, textStatus, jQxhr){
               parsedJSON = JSON.parse(response)
               currSecurity = parsedJSON["id"]
               currSymbol = parsedJSON["symbol"]
               var chLine = document.getElementById('priceLine');

               document.getElementById("securitySearch").value = "";
               currChart.destroy();

               GenerateNewPriceLine()
        },
        error: function( jqXhr, textStatus, errorThrown ){
                console.log( errorThrown );
        }

    })
}

function GetSymbolOfCurrSecurity(){
    var BearerString = 'Bearer ' + api.readCookie("accessToken")

    $.ajax({
            url: '/api/security/'+currSecurity,
            headers: {
                'Authorization': BearerString
            },
            contentType: 'application/json',
            type: "GET",
            success: function(response, textStatus, jQxhr){
                   parsedJSON = JSON.parse(response)
                   currSymbol = parsedJSON["symbol"]

            },
            error: function( jqXhr, textStatus, errorThrown ){
                    console.log( errorThrown );
            }

        })
}

function giveUserMoney(event) {
    amountMoneyGiven = parseFloat(document.getElementById("amountMoneyGiven").value);

    var BearerString = 'Bearer ' + api.readCookie("accessToken")

    giveMoneyParams = {
        "moneyIncrease" : amountMoneyGiven
    }

    $.ajax({
            url: '/api/customer/giveMoney',
            headers: {
                'Authorization': BearerString
            },
            contentType: 'application/json',
            type: "POST",
            data: JSON.stringify(giveMoneyParams),
            success: function(response, textStatus, jQxhr){
                   console.log(response)
                   document.location.reload()

            },
            error: function( jqXhr, textStatus, errorThrown ){
                    console.log( errorThrown );
            }

        });


}

function placeOrder(event) {
    numShares = parseInt(document.getElementById("numShares").value);
    amountPerShare = parseFloat(document.getElementById("amountPerShare").value);

    var BearerString = 'Bearer ' + api.readCookie("accessToken")


    orderParams = {
        "userId":currUser["id"],
        "investorAction":currAction,
        "investorType":0,
        "orderType":0,
        "symbol":currSymbol,
        "numShares": numShares,
        "costPerShare": amountPerShare,
        "timeCreated": moment().unix(),
        "allowTakers":false,
        "limitPerShare":0.0,
        "stopPrice":0.0
    }

    $.ajax({
        url: '/api/order',
        headers: {
            'Authorization': BearerString
        },
        contentType: 'application/json',
        type: "POST",
        data: JSON.stringify(orderParams),
        success: function(response, textStatus, jQxhr){
               console.log(response)
               document.location.reload()


        },
        error: function( jqXhr, textStatus, errorThrown ){
                console.log( errorThrown );
        }

    });


}
