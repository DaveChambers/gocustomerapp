function deleteCustomer(id, callback) {
	if (confirm("Are you sure you wish to remove this record?")) {
		fetch("http://localhost:8080/deletecustomer", {
			method: "POST",
			body: JSON.stringify(
				{id: id}
			),
			headers: {
				"Content-type": "application/json; charset=UTF-8"
			}
		}).then((response) => {
			callback();
		});
	}
}

function createNewCustomer() {
	window.location.href = "http://localhost:8080/create/";
}

function viewAllCustomers() {
    window.location.href = "http://localhost:8080/search/";
}

function editCustomer(id) {
    window.location.href = `http://localhost:8080/edit/${id}`;
}

function viewCustomer(id) {
    window.location.href = `http://localhost:8080/show/${id}`;
}