$("form#create-survey-form").on("submit", function(e) {
    e.preventDefault();

    var formData = {
        name: $("#survey-name").val(),
        team: {
            name: $("#team-name").val()
        },
        components: []
    }

    if ($("input[name=component-westrum]").prop("checked")) {
        formData.components.push({
            name: "westrum",
            category: "culture",
            weight: 1
        })
    }

    if ($("input[name=component-lencioni]").prop("checked")) {
        formData.components.push({
            name: "lencioni",
            category: "cohesion",
            weight: 2
        })
    }

    if ($("input[name=component-metrics]").prop("checked")) {
        formData.components.push({
            name: "metrics",
            category: "cohesion",
            weight: 3
        })
    }

    if ($("#team-new").val().length) {
        formData.team.name = $("#team-new").val()
    }

    $.ajax({
        type: "POST",
        url: '/api/v1/survey/',
        data: JSON.stringify(formData),
        contentType: "application/json; charset=utf-8",
        dataType: "json",
        success: function(data) {
            s = window.location.href.split("/")
            url = s[0] + "//" + s[2] + "/survey/" + formData.name
            dashUrl = `${s[0]}//${s[2]}/public/dashboard/survey/${formData.name}/?share_code=${data.share_code}`
            $("#form-success-link").text(url).prop('href', url)
            $("#dashboard-link").text(dashUrl).prop('href', dashUrl)
            $("#new-form").css('display', 'none')
            $("#form-success").css('display', 'block')
        },
        error: function(data) {
            if (data.responseText == '{"error":"UNIQUE constraint failed: surveys.name"}') {
                $("#request-error").html("A survey already exists with that name. Please try another name.").show();
            } else {
                $("#request-error").text(data)
            }
        }
    });
})