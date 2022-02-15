document.addEventListener('DOMContentLoaded', (event) => {

    document.querySelectorAll('.js-only').forEach((i) => {
        i.classList.toggle('js-only');
    });

    document.querySelectorAll('.cvss-preset').forEach((i) => {
        i.onclick = function (e) {
            document.getElementById('mincvss').value = this.getAttribute('data-preset');
        };
    });

    document.getElementById('date').value = today();

    // document.getElementById('composer-form').onsubmit = function (e) {
    //     document.getElementById('loading').classList.remove('d-none');
    // };

});

function today() {
    date = new Date();
    day = date.getDate();
    month = niceMonth(date.getMonth());
    year = date.getFullYear();
    return day + ' ' + month + ' ' + year;
}

function niceMonth(indx) {
    return [
        "January",
        "February",
        "March",
        "April",
        "May",
        "June",
        "July",
        "August",
        "September",
        "October",
        "November",
        "December",
    ][indx]
}