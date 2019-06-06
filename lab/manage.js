import {Spinner} from 'spin.js';

dict = []
versionOptions = ['2.15', '2.16', '2.17', '2.18', '2.19', '2.20', '2.21', '2.22', '2.23', '2.24', '2.25', '2.27']

$(document).ready(function() {
    console.log('Getting unit data')

    $.ajax({
        type: 'GET',
        url: '/api/units/',
        dataType: 'json',
        success: (data) => { 
            dict = data
            buildTable()
        }
    })
})

function buildTable() {
    // remove old table if there is one
    $('#dictTable tbody tr').remove()

    let table = $('#dictTable tbody')
/*
    let options
    for (let i = 0; i < versionOptions.length; i++) {
        if (versionOptions[i]) {
            options += $('<option value="' + i + '">').html(versionOptions[i])
        }
    }
*/
    for (let key in dict) {
        let curUnit = dict[key]

        let row = $('<tr>').appendTo(table)
        let nameElement = $('<td>').appendTo(row)
        let beanElement = $('<td>').appendTo(row)
        let versionElement = $('<td>').appendTo(row)
        let stateElement = $('<td>').appendTo(row)
        let dropdown = $('<select />').appendTo(versionElement)

        nameElement.append(
            $('<input>', {
                type: 'text',
                val: curUnit.name
            })
        )

        versionOptions.forEach(val => {
            $('<option />', {value: val, text: val}).appendTo(dropdown)
        })
        
        beanElement.append($('<div>').html(curUnit.beanID))
        stateElement.append(makeIcon(curUnit.state))
    }
}

function makeIcon(state) {
    switch (state) {
        case 0:
            //idle
            let check = $('<svg class="checkmark" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 52 52"> \
                <circle class="checkmark__circle" cx="26" cy="26" r="25" fill="none"/> \
                <path class="checkmark__check" fill="none" d="M14.1 27.2l7.1 7.2 16.7-16.8"/> </svg>')
            return check
        case 1:
            // updating
            let spinner = new Spinner().spin()
            return spinner.el
        case 2:
            // failed
            return $('<span class="ui-icon ui-icon-closethick">')
        default:
            console.log('unexpected state in units response: ' + state)
            return
    }
}