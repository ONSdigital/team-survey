var Dashboard = {
    data: {},
    team: "All",
    dateStart: "",
    dateEnd: "",
    survey: "",
    share_code: "",

    teams: {},
    init: function() {
        this.setListeners();
        this.getAllTeams();

        today = new Date()
        this.dateEnd = this.formatDate(today);
        this.dateStart = this.formatDate(new Date(today.getFullYear(), today.getMonth(), today.getDate() - 30));

        $("form.params input[name=dateStart]").val(this.dateStart)
        $("form.params input[name=dateEnd]").val(this.dateEnd)

        return this;
    },
    render: function() {
        this.renderLencioni();
        this.renderMetrics();
        this.renderWestrum();
        this.renderDownload();

        return this;
    },
    setListeners: function() {
        var d = this;
        $("form.params select").on('change', function(e) {
            var t = e.target.value;
            d.setTeam(t);
        });

        $("form.params input[name=dateStart]").on('change', function(e) {
            var t = e.target.value;
            d.dateStart = t;

            d.getTeamData()
            d.render();
        });

        $("form.params input[name=dateEnd]").on('change', function(e) {
            var t = e.target.value;
            d.dateEnd = t;

            d.getTeamData()
            d.render();
        });
    },
    setTeam: function(team) {
        this.team = team;
        this.getTeamData();
        this.render();
        $("#team-name").text("Team: " + team).prop('href', `/admin/dashboard/team/${team}/`)
        $("title").text("Dashboard - " + team);
        $("select[name=team]").val(team)
    },
    setSurvey: function(survey, share_code) {
        this.survey = survey;
        this.share_code = share_code;
        this.getSurveyData();
        this.render();
        $("#survey-name").text(survey)
        $("title").text("Dashboard - " + survey);
    },
    getSurveyData: async function() {
        var d = this;
        d.data = {};
        $.ajax({
            url: `/api/v1/survey/${d.survey}/stats/?share_code=${d.share_code}`,
            async: false,
            success: function(result) {
                d.data = result;
                $("#team-name").text(`(${d.data.team})`)
            }
        });

        return d;
    },
    getTeamData: async function() {
        var d = this;
        d.data = {};
        $.ajax({
            url: `/api/v1/stats/?team=${d.team}&dateStart=${d.dateStart}&dateEnd=${d.dateEnd}`,
            async: false,
            success: function(result) {
                d.data = result;
            },
            error: function(response) {
                alert(response.responseJSON.error)
            }
        });

        return d;
    },
    getAllTeams: async function() {
        var d = this;
        if ($(".sub-navbar .params select[name=team]").length) {
            $.ajax({
                url: `/api/v1/teams`,
                async: false,
                success: function(result) {
                    d.teams = result;
                    var ele = $(".sub-navbar .params select[name=team]").empty();
                    ele.append(`<option>All</option>`)
                    d.teams.forEach(function(t) {
                        ele.append(`<option>${t.name}</option>`)
                    })
                },
                error: function(response) {
                    alert(response.responseJSON.error)
                }
            });
        }
        return d;
    },
    renderDownload: function() {
        if (this.data.length == 0 || this.survey.length == 0) {
            $("#download-pack").hide();
            return false;
        }

        $("#download-pack").show();

        $("#download-pack").attr('href', "https://adoption-pack-generator-packs-prod.s3.amazonaws.com/" + encodeURIComponent(this.survey) + ".pptx")
    },
    renderLencioni: function() {
        if (this.data.length == 0) {
            console.log("No data");
            return false;
        }

        var d = this.data;

        $("#lencioni-trust").attr('class', 'trust ' + d.lencioni_trust)
        $("#lencioni-conflict").attr('class', 'conflict ' + d.lencioni_conflict)
        $("#lencioni-commitment").attr('class', 'commitment ' + d.lencioni_commitment)
        $("#lencioni-accountability").attr('class', 'accountability ' + d.lencioni_accountability)
        $("#lencioni-results").attr('class', 'results ' + d.lencioni_results)
    },
    renderWestrum: function() {
        var d = this.data;
        var p = ((this.data.westrum_score - 1) / 6) * 100

        console.log("PERCENTAGE: " + p)

        $(".westrum .westrum-scale .marker").css("left", p + "%")
    },
    renderMetrics: function() {
        var d = this.data;
        $('*[data-metric="mttr"]').text(d.mttr || "No Data");
        $('*[data-metric="lead_time"]').text(d.lead_time || "No Data");
        $('*[data-metric="deployment_frequency"]').text(d.deployment_frequency || "No Data");
        $('*[data-metric="change_failure"]').text(d.change_failure || "No Data");
        $('*[data-metric="westrum_result"]').text(d.westrum_result || "No Data");

        var responseText = `${d.total_responses || 0} responses`
        if (d.total_responses == 1) {
            responseText = `${d.total_responses} response`
        }
        $("#total-responses").text(responseText);
    },
    formatDate: function(date) {
        var day = date.getDate().padZero();
        var month = (date.getMonth() + 1).padZero();
        var year = date.getFullYear();

        return year + '-' + month + '-' + day;
    }
}

Number.prototype.padZero = function(len) {
    var s = String(this),
        c = '0';
    len = len || 2;
    while (s.length < len) s = c + s;
    return s;
}