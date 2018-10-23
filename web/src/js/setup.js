$(document).ready(function(){

})

function validateEmail(email) {
    var re = /^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
    return re.test(String(email).toLowerCase());
}

function registerUser(event){
    var FullName = document.getElementById("name").value;
    var Email = document.getElementById("email").value;
    var Password = document.getElementById("password").value;
    var PassCheck = document.getElementById("confirm").value;

    var fullNameArr = FullName.split(" ");
    if (fullNameArr.length != 2){
        alert("Please use your full name, in the following format: <First Name> <Last Name>");
        return
    }

    if (Password != PassCheck){
        alert("Password and Password confirm do not match!");
        return
        //Error
    }

    if (!validateEmail(Email)){
        alert("Email is of invalid structure");
        return
        //Error
    }

    userData = {
        "email": Email,
        "firstName": fullNameArr[0],
        "lastName": fullNameArr[1],
        "password": Password,
        "passwordVerify": PassCheck
    };

    console.log(userData);

    $.ajax({
        url: "/api/customer",
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


