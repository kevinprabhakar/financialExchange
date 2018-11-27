$(document).ready(function(){
      var BearerString = 'Bearer ' + api.readCookie("accessToken")
      $.ajax({
        url: "/api/customer/ownedShares",
        headers: {
            'Authorization': BearerString
        },
        type: "GET",
        contentType: "application/json",
        success: function(response, textStatus, jQxhr){
               ownedSharesList = JSON.parse(response)
               console.log(ownedSharesList)
               newHTMLString = ""
               for (var i = 0; i < ownedSharesList.length; i++){
                    newHTMLString += "<tr>"
                    newHTMLString += "<th scope='row'>"+(i+1)+"</th>"
                    newHTMLString += "<td>"+ownedSharesList[i].symbol+"</td>"
                    newHTMLString += "<td>"+ownedSharesList[i].entityName+"</td>"
                    newHTMLString += "<td>"+ownedSharesList[i].numShares+"</td>"
                    newHTMLString += "<td>$"+ownedSharesList[i].currPrice+"</td>"
                    newHTMLString += "</tr>"
               }
               var ownedSharesTable = document.getElementById("ownedSharesList")
               ownedSharesTable.innerHTML = newHTMLString

        },
        error: function( jqXhr, textStatus, errorThrown ){
                console.log( errorThrown );
        }
      })
});