describe('Team-Survey', function() {
    describe('I setup the survey', function() {
        it('it creates the survey', function() {
            cy.request({
                url: `${Cypress.env('API_SERVER')}/survey/`,
                method: "POST",
                failOnStatusCode: false,
                body: {
                    "name": Cypress.env('TEST_SURVEY_NAME'),
                    "team": {
                        "name": Cypress.env('TEST_SURVEY_TEAM')
                    }
                }
            })
        })
    })
    describe('I visit the first survey page', function() {
        beforeEach(function() {
            cy.visit(`${Cypress.env('HOST')}/survey/${Cypress.env('TEST_SURVEY_NAME')}`)
        })
        it('it has a page title and description', function() {
            cy.contains('Anonymous Team Questionnaire')
            cy.contains(`We're going to ask you a few questions about yourself, your team, and the software that your team looks after. All responses are anonymous, and won't be used against you for nefarious purposes.`)
        })
        it('it displays the correct questions in order', function() {
            cy.get('#surveyResult').should('not.be.visible')
            cy.get('#surveyElement').should('be.visible')

            cy.get('#surveyElement').contains('What is your current role?')
            cy.get('#surveyElement #sq_100 select')
                .should('contain', 'Choose...')
                .should('contain', 'Other')
            cy.get('#surveyElement #sq_101').should('not.be.visible')

            cy.get('#surveyElement').contains('What is your employment status?')
            cy.get('#surveyElement #sq_102 select')
                .should('contain', 'Choose...')
                .should('contain', 'Service Worker')

            cy.get('#surveyElement').contains('Which programming languages are you working with?')

            cy.get('#surveyElement #sq_103')
                .should('contain', 'Don\'t know')
                .should('contain', 'Other')

            cy.get('#surveyElement #sq_104').should('not.be.visible')

            cy.get('#surveyElement').get('input[type=button].sv_prev_btn').should('not.be.visible').should('have.value', 'Previous')
            cy.get('#surveyElement').get('input[type=button].sv_next_btn').should('be.visible').should('have.value', 'Next')
        })
        it('it allows input for other options', function() {
            cy.get('#surveyElement #sq_101').should('not.be.visible')
            cy.get('#surveyElement #sq_100 select').select('Other')
            cy.get('#surveyElement #sq_101').should('be.visible')
            cy.get('#surveyElement #sq_100 select').select('Choose...')
            cy.get('#surveyElement #sq_101').should('not.be.visible')

            cy.get('#surveyElement #sq_104').should('not.be.visible')
            cy.get('#surveyElement #sq_103').contains('Other').click()
            cy.get('#surveyElement #sq_104').should('be.visible')
            cy.get('#surveyElement #sq_103').contains('Other').click()
            cy.get('#surveyElement #sq_104').should('not.be.visible')
        })
        it('it validates correctly', function() {
            cy.get('.sv_q_erbox').should('not.be.visible')
            cy.get('#surveyElement').get('input[type=button].sv_next_btn').click()
            cy.get('#surveyElement #sq_100 .sv_q_erbox').should('be.visible').should('contain', 'Please answer the question')
            cy.get('#surveyElement #sq_100 select').select('Other')
            cy.get('#surveyElement').get('input[type=button].sv_next_btn').click()
            cy.get('#surveyElement #sq_100 .sv_q_erbox').should('not.be.visible')

            cy.get('#surveyElement #sq_101 .sv_q_erbox').should('be.visible').should('contain', 'Please answer the question')
            cy.get('#surveyElement #sq_100 select').select('User Researcher')
            cy.get('#surveyElement').get('input[type=button].sv_next_btn').click()
            cy.get('#surveyElement #sq_100 .sv_q_erbox').should('not.be.visible')

            cy.get('#surveyElement #sq_100 select').select('Other')
            cy.get('#surveyElement').get('input[type=button].sv_next_btn').click()
            cy.get('#surveyElement #sq_100 .sv_q_erbox').should('not.be.visible')
            cy.get('#surveyElement #sq_101 .sv_q_erbox').should('be.visible')

            cy.get('#surveyElement #sq_101 input.sv_q_text_root').should('be.visible').type('Tester')
            cy.get('#surveyElement').get('input[type=button].sv_next_btn').click()
            cy.get('#surveyElement #sq_100 .sv_q_erbox').should('not.be.visible')
            cy.get('#surveyElement #sq_101 .sv_q_erbox').should('not.be.visible')

            cy.get('#surveyElement #sq_102 .sv_q_erbox').should('be.visible')
            cy.get('#surveyElement #sq_102 select').select('Contractor')
            cy.get('#surveyElement').get('input[type=button].sv_next_btn').click()
            cy.get('#surveyElement #sq_102 .sv_q_erbox').should('not.be.visible')

            cy.get('#surveyElement #sq_103 .sv_q_erbox').should('be.visible')
            cy.get('#surveyElement #sq_103').contains('Golang').click()
            cy.get('#surveyElement #sq_103').contains('Other').click()

            cy.get('#surveyElement').get('input[type=button].sv_next_btn').click()

            cy.get('#surveyElement #sq_104 .sv_q_erbox').should('be.visible')
            cy.get('#surveyElement #sq_104 input.sv_q_text_root').should('be.visible').type('Assembly')

            cy.get('#surveyElement').get('input[type=button].sv_next_btn').click()
            cy.get('.sv_q_erbox').should('not.be.visible')

            cy.get('#surveyElement #sq_100').should('not.be.visible')
            cy.get('#surveyElement #sp_101').should('be.visible')
        })
    })

    describe('I visit the second survey page', function() {
        it('it has a page title and description', function() {
            cy.contains('Anonymous Team Questionnaire')
            cy.contains(`It is important to evaluate the statements honestly and without over-thinking your answers. Remember, this is anonymous and you won't be tickled for being truthful.`)
            cy.contains('Thinking about your team, how frequently do the following statements apply?')
        })
        it('it displays the correct questions in order', function() {
            cy.contains('Strongly Disagree')
            cy.contains('Disagree')
            cy.contains('Neither Agree nor Disagree')
            cy.contains('Agree')
            cy.contains('Strongly Agree')

            cy.contains(`Information is actively sought.`)
            cy.contains(`Messengers are not punished when they deliver news of failures or other bad news.`)
            cy.contains(`Responsibilities are shared.`)
            cy.contains(`Cross-functional collaboration is encouraged or rewarded.`)
            cy.contains(`Failure causes enquiry.`)
            cy.contains(`New ideas are welcomed.`)
            cy.contains(`Failures are treated primarily as opportunities to improve the system.`)

            cy.get('#surveyElement').get('input[type=button].sv_prev_btn').should('be.visible').should('have.value', 'Previous')
            cy.get('#surveyElement').get('input[type=button].sv_next_btn').should('be.visible').should('have.value', 'Next')
        })

        it('it validates correctly', function() {
            cy.get('.sv_q_erbox').should('not.be.visible')
            cy.get('#surveyElement').get('input[type=button].sv_next_btn').click()
            cy.get('.sv_q_erbox').should('be.visible')

            cy.get('#surveyElement table.sv_q_matrix tbody tr').each(function(e, i, l) {
                cy.wrap(l[i]).find('input[type=radio]').eq(Math.floor((Math.random() * 5))).click()
            })
            cy.get('#surveyElement').get('input[type=button].sv_next_btn').click()
            cy.get('.sv_q_erbox').should('not.be.visible')

            cy.get('#surveyElement #sp_101').should('not.be.visible')
            cy.get('#surveyElement #sp_102').should('be.visible')
        })
    })

    describe('I visit the third survey page', function() {
        it('it has a page title and description', function() {
            cy.contains('Anonymous Team Questionnaire')
            cy.contains(`It is important to evaluate the statements honestly and without over-thinking your answers. Remember, this is anonymous and we won't harm your pets if we don't like your answers.`)
            cy.contains('Thinking about your team, do you agree with the following statements?')
        })
        it('it displays the correct questions in order', function() {
            cy.contains('Rarely')
            cy.contains('Sometimes')
            cy.contains('Usually')

            cy.contains(`Team members are passionate and unguarded in their discussion of issues.`)
            cy.contains(`Team members call out one another’s deficiencies or unproductive behaviours.`)
            cy.contains(`Team members know what their peers are working on and how they contribute to the collective good of the team.`)
            cy.contains(`Team members quickly and genuinely apologise to one another when they say or do something inappropriate or possibly damaging to the team.`)
            cy.contains(`Team members willingly make sacrifices (such as budget, turf, head count) in their departments or areas of expertise for the good of the team.`)
            cy.contains(`Team members openly admit their weaknesses and mistakes.`)
            cy.contains(`Team meetings are compelling, not boring.`)
            cy.contains(`Team members leave meetings confident that their peers are completely committed to the decisions that were agreed on, even if they were in initial disagreement.`)
            cy.contains(`Morale is significantly affected by the failure to achieve team goals.`)
            cy.contains(`During team meetings, the most important—and difficult—issues are put on the table to be resolved.`)
            cy.contains(`Team members are deeply concerned about the prospect of letting down their peers.`)
            cy.contains(`Team members know about one another’s personal lives and are comfortable discussing them.`)
            cy.contains(`Team members end discussions with clear and specific resolutions and action plans.`)
            cy.contains(`Team members challenge one another about their plans and approaches.`)
            cy.contains(`Team members are slow to seek credit for their own contributions, but quick to point out those of others.`)


            cy.get('#surveyElement').get('input[type=button].sv_prev_btn').should('be.visible').should('have.value', 'Previous')
            cy.get('#surveyElement').get('input[type=button].sv_next_btn').should('be.visible').should('have.value', 'Next')
        })
        it('it validates correctly', function() {
            cy.get('.sv_q_erbox').should('not.be.visible')
            cy.get('#surveyElement').get('input[type=button].sv_next_btn').click()
            cy.get('.sv_q_erbox').should('be.visible')

            cy.get('#surveyElement table.sv_q_matrix tbody tr').each(function(e, i, l) {
                cy.wrap(l[i]).find('input[type=radio]').eq(Math.floor((Math.random() * 3))).click()
            })
            cy.get('#surveyElement').get('input[type=button].sv_next_btn').click()
            cy.get('.sv_q_erbox').should('not.be.visible')

            cy.get('#surveyElement #sp_102').should('not.be.visible')
        })
    })
    describe('I visit the fourth survey page', function() {
        it('it has a page title and description', function() {
            cy.contains('Anonymous Team Questionnaire')
            cy.contains(`We're asking all the teams to measure at least four key metrics, they are Lead Time, Deployment Frequency, Mean Time to Recovery (MTTR) and Change Failure Percentage.`)
        })

        it('it displays the correct questions in order', function() {
            cy.contains(`For the primary application or service you work on, what is your lead time for changes (ie. how long does it take to go from code committed to code successfully running in production)? `)
            cy.contains(`For the primary application or service you work on, how often does your organisation deploy code into production or release it to end users?`)
            cy.contains(`For the primary application or service you work on, how long does it generally take to restore service when an incident or a defect that impacts users occurs (e.g. unplanned outage, service impairment)?`)
            cy.contains(`For the primary application or service you work on, what percentage of changes to production or released to users result in degraded service (e.g., lead to service impairment or service outage) and subsequently require remediation (e.g. require a hotfix, a rollback, a fix-forward, or a patch)? `)

            cy.get('#surveyElement').get('input[type=button].sv_prev_btn').should('be.visible').should('have.value', 'Previous')
            cy.get('#surveyElement').get('input[type=button].sv_complete_btn').should('be.visible').should('have.value', 'Complete')
        })
        it('it validates correctly', function() {
            cy.get('.sv_q_erbox').should('not.be.visible')
            cy.get('#surveyElement').get('input[type=button].sv_complete_btn').click()

            cy.get('#sq_107 .sv_q_erbox').should('be.visible')
            cy.get('#sq_107 select').select('Less than one day')
            cy.get('#surveyElement').get('input[type=button].sv_complete_btn').click()
            cy.get('#sq_107 .sv_q_erbox').should('not.be.visible')

            cy.get('#sq_108 .sv_q_erbox').should('be.visible')
            cy.get('#sq_108 select').select('Don\'t know')
            cy.get('#surveyElement').get('input[type=button].sv_complete_btn').click()
            cy.get('#sq_108 .sv_q_erbox').should('not.be.visible')

            cy.get('#sq_109 .sv_q_erbox').should('be.visible')
            cy.get('#sq_109 select').select('Less than one day')
            cy.get('#surveyElement').get('input[type=button].sv_complete_btn').click()
            cy.get('#sq_109 .sv_q_erbox').should('not.be.visible')

            cy.get('#sq_110 .sv_q_erbox').should('be.visible')
            cy.get('#sq_110 select').select('16%-30%')
            cy.get('#surveyElement').get('input[type=button].sv_complete_btn').click()
            cy.get('#sq_110 .sv_q_erbox').should('not.be.visible')
        })
    })

    describe('I have completed the survey', function() {
        it('it has a page title and message', function() {
            cy.contains('Anonymous Team Questionnaire')
            cy.contains(`Thank you for your feedback.`)
            cy.contains(`All responses are anonymous, and won't be used against you for nefarious purposes!`)
        })
    })
})