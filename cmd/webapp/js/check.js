document.querySelectorAll(".row-checkbox").forEach(function (checkbox) {
  checkbox.addEventListener("change", function () {
    if (!this.checked) {
      this.checked = true;
    }

    const id = this.getAttribute("data-id");

    fetch("/tasks/done/" + id, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ checked: this.checked }),
    })
      .then((response) => {
        if (!response.ok) {
          console.error("Request failed.");
        } else {
          const statusCell =
            this.closest("tr").querySelector("td:nth-child(3)");
          if (statusCell) {
            statusCell.textContent = "done";
          }
        }
      })
      .catch((error) => {
        console.error("Network error:", error);
      });
  });
});
