function uploadFile() {
	var xhr = new XMLHttpRequest();
	xhr.open('POST', 'localhost:1112/storage', true);
	xhr.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
	xhr.onload = function () {
	    // do something to response
	    console.log(this.responseText);
	};
	xhr.send('user=person&pwd=password&organization=place&requiredkey=key');
}