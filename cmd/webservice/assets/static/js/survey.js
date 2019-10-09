var json = {
    "completedHtml": "<h3>Thank you for your feedback.</h3> <h5>All responses are anonymous, and won't be used against you for nefarious purposes!</h5>",
    "pages": [{
            "name": "page1",
            "title": "We're going to ask you a few questions about yourself, your team, and the software that your team looks after. All responses are anonymous, and won't be used against you for nefarious purposes.",
            "elements": [{
                    type: "dropdown",
                    name: "current_role",
                    title: "What is your current role?",
                    isRequired: true,
                    colCount: 0,
                    choices: [
                        "User Researcher",
                        "Business Relationship Manager",
                        "Senior Delivery Manager",
                        "Delivery Manager",
                        "Service Manager (non digital)",
                        "Service Manager (digital)",
                        "Agile Coach",
                        "Embedded SME",
                        "SLT",
                        "HR/Recruitment",
                        "Product Owner",
                        "QA Consultant",
                        "Business Analyst",
                        "Interaction Designer",
                        "Solution Architect",
                        "Technical Architect",
                        "Technical Lead",
                        "Engineer - SEO",
                        "Engineer - HEO",
                        "Engineer - EO",
                        "Other"
                    ]
                },
                {
                    name: "current_role_other",
                    title: "What is your current role?",
                    visibleIf: "{current_role}='Other'",
                    isRequired: true,
                    type: "text"
                },
                {
                    type: "dropdown",
                    name: "employment_status",
                    title: "What is your employment status?",
                    isRequired: true,
                    colCount: 0,
                    choices: [
                        "Permanent Employee",
                        "Contractor",
                        "Service Worker",
                    ]
                },
                {
                    type: "checkbox",
                    name: "programming_languages",
                    title: "Which programming languages are you working with?",
                    isRequired: true,
                    colCount: 0,
                    choices: [
                        "C",
                        "C++",
                        "Golang",
                        "Java",
                        "Node/Javascript",
                        "Python",
                        "Ruby",
                        "Scala",
                        "Other",
                        "Don't know"
                    ]
                },
                {
                    type: "text",
                    name: "programming_languages_other",
                    title: "Other programming languages: ",
                    visibleIf: "{programming_languages} contains 'Other'",
                    isRequired: true
                },
            ]
        },
        {
            "name": "page2",
            "title": "It is important to evaluate the statements honestly and without over-thinking your answers. Remember, this is anonymous and you won't be tickled for being truthful.",
            "elements": [{
                type: "matrix",
                name: "about_the_team",
                title: "Thinking about your team, how frequently do the following statements apply?",
                isRequired: true,
                columns: [{
                    value: 1,
                    text: "Strongly Disagree"
                }, {
                    value: 2,
                    text: "Disagree"
                }, {
                    value: 3,
                    text: "Neutral"
                }, {
                    value: 4,
                    text: "Agree"
                }, {
                    value: 5,
                    text: "Strongly Agree"
                }],
                rows: [{
                    value: "information_actively_sought",
                    text: "Information is actively sought."
                }, {
                    value: "messengers_not_punished",
                    text: "Messengers are not punished when they deliver news of failures or other bad news."
                }, {
                    value: "responsibilities_shared",
                    text: "Responsibilities are shared."
                }, {
                    value: "collaboration_encouraged",
                    text: "Cross-functional collaboration is encouraged or rewarded."
                }, {
                    value: "failure_causes_enquiry",
                    text: "Failure causes enquiry."
                }, {
                    value: "new_ideas_welcomed",
                    text: "New ideas are welcomed."
                }, {
                    value: "failure_treated_as_opportunity",
                    text: "Failures are treated primarily as opportunities to improve the system."
                }]
            }]
        },
        {
            "name": "page3",
            "title": "It is important to evaluate the statements honestly and without over-thinking your answers. Remember, this is anonymous and we won't harm your pets if we don't like your answers.",
            "elements": [{
                type: "matrix",
                name: "about_the_team_extended",
                isRequired: true,
                title: "Thinking about your team, do you agree with the following statements?",
                columns: [{
                    value: 1,
                    text: "Rarely"
                }, {
                    value: 2,
                    text: "Sometimes"
                }, {
                    value: 3,
                    text: "Usually"
                }, ],
                rows: [{
                    value: "unguarded_discussion",
                    text: "Team members are passionate and unguarded in their discussion of issues."
                }, {
                    value: "call_out_unproductive_behaviour",
                    text: "Team members call out one another’s deficiencies or unproductive behaviours."
                }, {
                    value: "contribute_to_collective_good",
                    text: "Team members know what their peers are working on and how they contribute to the collective good of the team."
                }, {
                    value: "apologise",
                    text: "Team members quickly and genuinely apologise to one another when they say or do something inappropriate or possibly damaging to the team."
                }, {
                    value: "willingly_make_sacrifices",
                    text: "Team members willingly make sacrifices (such as budget, turf, head count) in their departments or areas of expertise for the good of the team."
                }, {
                    value: "admit_mistakes",
                    text: "Team members openly admit their weaknesses and mistakes."
                }, {
                    value: "compelling_meetings",
                    text: "Team meetings are compelling, not boring."
                }, {
                    value: "leave_meetings_committed",
                    text: "Team members leave meetings confident that their peers are completely committed to the decisions that were agreed on, even if they were in initial disagreement."
                }, {
                    value: "morale_affected_by_failure",
                    text: "Morale is significantly affected by the failure to achieve team goals."
                }, {
                    value: "difficult_issues_raised",
                    text: "During team meetings, the most important—and difficult—issues are put on the table to be resolved."
                }, {
                    value: "concerned_about_letting_down_peers",
                    text: "Team members are deeply concerned about the prospect of letting down their peers."
                }, {
                    value: "comfortable_discussing_personal_lives",
                    text: "Team members know about one another’s personal lives and are comfortable discussing them."
                }, {
                    value: "clear_resolution_discussions",
                    text: "Team members end discussions with clear and specific resolutions and action plans."
                }, {
                    value: "challenge_one_another",
                    text: "Team members challenge one another about their plans and approaches."
                }, {
                    value: "slow_to_seek_credit",
                    text: "Team members are slow to seek credit for their own contributions, but quick to point out those of others."
                }]
            }]
        },
        {
            "name": "page4",
            "title": "We're asking all the teams to measure at least four key metrics, they are Lead Time, Deployment Frequency, Mean Time to Recovery (MTTR) and Change Failure Percentage.",
            isRequired: true,
            "elements": [{
                    type: "dropdown",
                    name: "lead_time",
                    title: "For the primary application or service you work on, what is your lead time for changes (ie. how long does it take to go from code committed to code successfully running in production)?",
                    isRequired: true,
                    colCount: 0,
                    choices: [
                        "Between one day and one week",
                        "Between one month and six months",
                        "Between one week and one month",
                        "Less than one day",
                        "Less than one hour",
                        "More than six months",
                        "Don't know"
                    ]
                },
                {
                    type: "dropdown",
                    name: "deployment_frequency",
                    title: "For the primary application or service you work on, how often does your organisation deploy code into production or release it to end users?",
                    isRequired: true,
                    colCount: 0,
                    choices: [
                        "Between once per day and once per week",
                        "Between once per hour and once per day",
                        "Between once per month and once every six months",
                        "Between once per week and once per month",
                        "Fewer than once every six months",
                        "On demand (multiple deploys per day)",
                        "Don't know"
                    ]
                },
                {
                    type: "dropdown",
                    name: "mttr",
                    title: "For the primary application or service you work on, how long does it generally take to restore service when an incident or a defect that impacts users occurs (e.g. unplanned outage, service impairment)?",
                    isRequired: true,
                    colCount: 0,
                    choices: [
                        "More than six months",
                        "Between one month and six months",
                        "Between one week and one month",
                        "Between one day and one week",
                        "Less than one day",
                        "Less than one hour",
                        "Don't know"
                    ]
                },
                {
                    type: "dropdown",
                    name: "change_failure",
                    title: "For the primary application or service you work on, what percentage of changes to production or released to users result in degraded service (e.g., lead to service impairment or service outage) and subsequently require remediation (e.g. require a hotfix, a rollback, a fix-forward, or a patch)?",
                    isRequired: true,
                    colCount: 0,
                    choices: [
                        "0%-15%",
                        "16%-30%",
                        "31%-45%",
                        "46%-60%",
                        "61%-75%",
                        "76%-100%",
                        "Don't know"
                    ]
                },
            ]
        },
    ],
    "showQuestionNumbers": "off"
};

window.survey = new Survey.Model(json);

survey
    .onComplete
    .add(function(result) {
        result.data["team"] = window.TeamName;
        var scriptUrl = "/api/v1/survey/";
        var r = flatten(result.data);
        r.team = window.TeamName;
        var bodyJSON = {
            name: window.SurveyName,
            id: parseInt(window.SurveyID),
            team: {
                name: window.TeamName
            },
            results: [r]
        };

        (async() => {
            const rawResponse = await fetch(scriptUrl, {
                method: 'POST',
                body: JSON.stringify(bodyJSON)
            });
            const content = await rawResponse.json();
        })();
    });

$("#surveyElement").Survey({ model: survey });

function flatten(obj, prefix, current) {
    prefix = prefix || []
    current = current || {}

    // Remember kids, null is also an object!
    if (typeof(obj) === 'object' && obj !== null) {
        Object.keys(obj).forEach(key => {
            this.flatten(obj[key], prefix.concat(key), current)
        })
    } else {
        current[prefix.join('.')] = obj
    }

    return current
}