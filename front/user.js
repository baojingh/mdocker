document.getElementById("loginForm").addEventListener("submit", function(event) {
    event.preventDefault(); // 阻止表单提交的默认行为

    var xhr = new XMLHttpRequest();
    var url = "http://127.0.0.1:8081/login";
    var username = document.getElementById("username").value;
    var password = document.getElementById("password").value;
    var data = JSON.stringify({ username: username, password: password });

    xhr.open("POST", url, true);
    xhr.setRequestHeader("Content-Type", "application/json");

    xhr.onreadystatechange = function() {
        if (xhr.readyState === XMLHttpRequest.DONE) {
            if (xhr.status === 200) {
                // 登录成功
                var response = JSON.parse(xhr.responseText);
                console.log(response);
            } else {
                // 登录失败
                console.error("登录失败");
            }
        }
    };

    xhr.send(data);
});