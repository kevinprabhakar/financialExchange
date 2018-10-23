$(document).ready(function(){

})

function validateEmail(email) {
    var re = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
    return re.test(String(email).toLowerCase());
}

function loginUser(event){
    var Email = document.getElementById("email").value;
    var Password = document.getElementById("password").value;

    if (!validateEmail(Email)){
        alert("Email is of invalid structure");
        return
        //Error
    }

    userData = {
        "email": Email,
        "password": Password,
    };


    $.ajax({
        url: "/api/customer/login",
        type: "POST",
        data: JSON.stringify(userData),
        success: function(response, textStatus, jQxhr){
            accessToken = JSON.parse(response);
            api.createCookie("accessToken", accessToken["accessToken"],1);
            window.location.replace("/prices.html");
        },
        error: function( jqXhr, textStatus, errorThrown ){
            console.log( errorThrown );
        }
    })
}


