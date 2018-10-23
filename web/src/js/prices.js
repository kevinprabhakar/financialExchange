var currAction = 0
var currUser = {}
var currSecurity = 1
var currTimePeriod = 0

$(document).ready(function(){
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
                  label: "prices",
                  data: pricePoints,
                }]
              };

              // Create a new line chart object where as first parameter we pass in a selector
              // that is resolving to our chart container element. The Second parameter
              // is the actual data object.
              var chLine = document.getElementById("priceLine");
              if (chLine) {
                new Chart(chLine, {
                type: 'line',
                data: chartData,
                });
              }
          },
          error: function( jqXhr, textStatus, errorThrown ){
                  console.log( errorThrown );
          }
      });
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
               console.log(currUser)
        },
        error: function( jqXhr, textStatus, errorThrown ){
                console.log( errorThrown );
        }
      })
});

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
    securityName = document.getElementById("securitySearch").value;

}

function placeOrder(event) {
    numShares = parseInt(document.getElementById("numShares").value);
    amountPerShare = parseFloat(document.getElementById("amountPerShare").value);


}
