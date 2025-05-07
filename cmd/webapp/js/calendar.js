const form = document.querySelector("form");
const input = document.getElementById("dueDate");
const calendarContainer = document.getElementById("calendarContainer");
const cal = document.getElementById("callyPicker");
const dueInput = document.getElementById("due");

input.addEventListener("click", () => {
  calendarContainer.classList.remove("hidden");
});

cal.addEventListener("change", (e) => {
  const selected = e.target.value;
  if (selected) {
    const date = new Date(selected);

    const displayFormatted = date.toDateString();
    input.value = displayFormatted;

    formatDateForGo(date);

    calendarContainer.classList.add("hidden");
  }
});

document.addEventListener("click", (e) => {
  if (!calendarContainer.contains(e.target) && e.target !== input) {
    calendarContainer.classList.add("hidden");
  }
});

if (!input.value) {
  const today = new Date();
  input.value = today.toDateString();
  formatDateForGo(today);
} else {
  try {
    const existingDate = new Date(input.value);
    formatDateForGo(existingDate);
  } catch (err) {
    console.error("Error formatting existing date:", err);
  }
}

function formatDateForGo(date) {
  const daysOfWeek = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];
  const months = [
    "Jan",
    "Feb",
    "Mar",
    "Apr",
    "May",
    "Jun",
    "Jul",
    "Aug",
    "Sep",
    "Oct",
    "Nov",
    "Dec",
  ];

  const dayOfWeek = daysOfWeek[date.getDay()];
  const month = months[date.getMonth()];
  const day = date.getDate();
  const year = date.getFullYear();

  const formattedDate = `${dayOfWeek} ${month} ${day} ${year}`;
  dueInput.value = formattedDate;
  console.log("Formatted date for Go:", formattedDate);
}

form.addEventListener("submit", function (e) {
  if (input.value) {
    try {
      const selectedDate = new Date(input.value);
      formatDateForGo(selectedDate);
    } catch (err) {
      e.preventDefault();
      alert("Please select a valid date");
    }
  } else {
    e.preventDefault();
    alert("Please select a due date");
  }
});
