describe('Team-Survey', function() {
    describe('I setup Team-Survey', function() {
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
    describe('I visit the homepage', function() {
        beforeEach(function() {
            cy.visit(`${Cypress.env('HOST')}`, {
                auth: {
                    username: Cypress.env('ADMIN_USERNAME'),
                    password: Cypress.env('ADMIN_PASSWORD')
                }
            })
        })
        it('it has a log in button and logo', function() {
            cy.contains('Log In')
        })
        describe('I click the log in button', function() {
            it('it takes me to the dashboard', function() {
                cy.contains('Log In').should('be.visible').click()
                cy.url().should('include', '/admin/dashboard/')
            })
        })
    })
    describe('I visit /admin/dashboard', function() {
        beforeEach(function() {
            cy.visit(`${Cypress.env('HOST')}/admin/dashboard/`, {
                auth: {
                    username: Cypress.env('ADMIN_USERNAME'),
                    password: Cypress.env('ADMIN_PASSWORD')
                }
            })
        })

        it('it displays results panels', function() {
            cy.contains('Westrum').should('be.visible')
            cy.get('.westrum').contains("Pathological")
            cy.get('.westrum').contains("Bureaucratic")
            cy.get('.westrum').contains("Generative")
            cy.contains('Lencioni').should('be.visible')
            cy.get('.pyramid').contains("Results")
            cy.get('.pyramid').contains("Accountability")
            cy.get('.pyramid').contains("Commitment")
            cy.get('.pyramid').contains("Conflict")
            cy.get('.pyramid').contains("Trust")
            cy.contains('Key Metrics').should('be.visible')
            cy.get('.metric-panel .pure-table').contains("Lead Time")
            cy.get('.metric-panel .pure-table').contains("Change Failure")
            cy.get('.metric-panel .pure-table').contains("MTTR")
            cy.get('.metric-panel .pure-table').contains("Deploy Frequency")
        })
        it('it displays the team name linked to the team dashboard', function() {
            cy.get('#team-name').should('be.visible')
                .should('have.text', 'Team: All')
                .should('have.prop', 'href', `${Cypress.env('HOST')}/admin/dashboard/team/All/`)
                .should('not.have.prop', 'target', '_blank')
        })
        it('it displays site navigation', function() {
            cy.get('.navbar').should('be.visible')
            cy.get('.navbar').contains("Dashboard").should('have.class', 'active')
            cy.get('.navbar').contains("Surveys").should('not.have.class', 'active')
        })
        it('it displays sub-navigation', function() {
            cy.get('.sub-navbar').should('be.visible')

            cy.get('.sub-navbar form').should('be.visible').should('have.class', 'params')
            cy.get('.sub-navbar form select[name=team]').contains('All')
            cy.get('.sub-navbar form').contains('From')
            cy.get('.sub-navbar form input[name=dateStart]').should('be.visible')
                .should('have.value', Cypress.moment().subtract(30, 'days').format('YYYY-MM-DD'))
            cy.get('.sub-navbar form').contains('To')
            cy.get('.sub-navbar form input[name=dateEnd]').should('be.visible')
                .should('have.value', Cypress.moment().format('YYYY-MM-DD'))
        })
        describe('I click the team dashboard link', function() {
            it('it takes me to the team dashboard', function() {
                cy.get('#team-name').click()
                cy.url().should('include', '/admin/dashboard/team/All/')
            })
        })
        describe('I click on the dashboard navigation link', function() {
            it('it takes me to the team dashboard', function() {
                cy.get('.navbar').contains('Dashboard').click()
                cy.url().should('equal', `${Cypress.env('HOST')}/admin/dashboard/`)
            })
        })
        describe('I click on the surveys navigation link', function() {
            it('it takes me to the survey list', function() {
                cy.get('.navbar').contains('Surveys').click()
                cy.url().should('equal', `${Cypress.env('HOST')}/admin/survey/`)
            })
        })
    });
    describe('I visit /admin/dashboard/team/All', function() {
        beforeEach(function() {
            cy.visit(`${Cypress.env('HOST')}/admin/dashboard/team/All/`, {
                auth: {
                    username: Cypress.env('ADMIN_USERNAME'),
                    password: Cypress.env('ADMIN_PASSWORD')
                }
            })
        })
        it('it displays results panels', function() {
            cy.contains('Westrum').should('be.visible')
            cy.get('.westrum').contains("Pathological")
            cy.get('.westrum').contains("Bureaucratic")
            cy.get('.westrum').contains("Generative")
            cy.contains('Lencioni').should('be.visible')
            cy.get('.pyramid').contains("Results")
            cy.get('.pyramid').contains("Accountability")
            cy.get('.pyramid').contains("Commitment")
            cy.get('.pyramid').contains("Conflict")
            cy.get('.pyramid').contains("Trust")
            cy.contains('Key Metrics').should('be.visible')
            cy.get('.metric-panel .pure-table').contains("Lead Time")
            cy.get('.metric-panel .pure-table').contains("Change Failure")
            cy.get('.metric-panel .pure-table').contains("MTTR")
            cy.get('.metric-panel .pure-table').contains("Deploy Frequency")

            cy.get('#team-name').should('be.visible')
                .should('have.text', 'Team: All')
                .should('have.prop', 'href', `${Cypress.env('HOST')}/admin/dashboard/team/All/`)
                .should('not.have.prop', 'target', '_blank')

            cy.get('h1.logo').should('be.visible')
                .should('have.text', 'Team-Survey')
        })
        it('it displays site navigation', function() {
            cy.get('.navbar').should('be.visible')
            cy.get('.navbar').contains("Dashboard").should('have.class', 'active')
            cy.get('.navbar').contains("Surveys").should('not.have.class', 'active')
        })
        it('it displays sub-navigation', function() {
            cy.get('.sub-navbar').should('be.visible')
            cy.get('.sub-navbar form').should('be.visible').should('have.class', 'params')
            cy.get('.sub-navbar form select[name=team]').should('not.exist')
            cy.get('.sub-navbar form').contains('From')
            cy.get('.sub-navbar form input[name=dateStart]').should('be.visible')
                .should('have.value', Cypress.moment().subtract(30, 'days').format('YYYY-MM-DD'))
            cy.get('.sub-navbar form').contains('To')
            cy.get('.sub-navbar form input[name=dateEnd]').should('be.visible')
                .should('have.value', Cypress.moment().format('YYYY-MM-DD'))
        })
        describe('I click on the dashboard navigation link', function() {
            it('it takes me to the team dashboard', function() {
                cy.get('.navbar').contains('Dashboard').click()
                cy.url().should('equal', `${Cypress.env('HOST')}/admin/dashboard/`)
            })
        })
        describe('I click on the surveys navigation link', function() {
            it('it takes me to the survey list', function() {
                cy.get('.navbar').contains('Surveys').click()
                cy.url().should('equal', `${Cypress.env('HOST')}/admin/survey/`)
            })
        })
    })
    describe('I visit /admin/survey', function() {
        beforeEach(function() {
            cy.visit(`${Cypress.env('HOST')}/admin/survey/`, {
                auth: {
                    username: Cypress.env('ADMIN_USERNAME'),
                    password: Cypress.env('ADMIN_PASSWORD')
                }
            })
        })
        it('it displays site navigation', function() {
            cy.get('.navbar').should('be.visible')
            cy.get('.navbar').contains("Dashboard").should('not.have.class', 'active')
            cy.get('.navbar').contains("Surveys").should('have.class', 'active')
        })
        it('it displays sub-navigation', function() {
            cy.get('.sub-navbar').should('be.visible')
            cy.get('.sub-navbar form').should('not.exist')
            cy.get('.sub-navbar').contains('All Surveys').should('have.class', 'active')
            cy.get('.sub-navbar').contains('Create New').should('not.have.class', 'active')
        })
        it('it displays a survey table', function() {
            cy.get('.page-panel table.pure-table').should('be.visible')
            cy.get('.page-panel table.pure-table thead').contains('Survey')
            cy.get('.page-panel table.pure-table thead').contains('# Responses')
            cy.get('.page-panel table.pure-table thead').contains('Created Date')
            cy.get('.page-panel table.pure-table thead').contains('Created Time')
        })
        describe('I click on the All Surveys sub-navigation link', function() {
            it('it takes me to the the survey list', function() {
                cy.get('.sub-navbar').contains('All Surveys').click()
                cy.url().should('equal', `${Cypress.env('HOST')}/admin/survey/`)
            })
        })
        describe('I click on the Create New sub-navigation link', function() {
            it('it takes me to the the survey list', function() {
                cy.get('.sub-navbar').contains('Create New').click()
                cy.url().should('equal', `${Cypress.env('HOST')}/admin/survey/new/`)
            })
        })
        describe('I click on the dashboard navigation link', function() {
            it('it takes me to the team dashboard', function() {
                cy.get('.navbar').contains('Dashboard').click()
                cy.url().should('equal', `${Cypress.env('HOST')}/admin/dashboard/`)
            })
        })
        describe('I click on the surveys navigation link', function() {
            it('it takes me to the survey list', function() {
                cy.get('.navbar').contains('Surveys').click()
                cy.url().should('equal', `${Cypress.env('HOST')}/admin/survey/`)
            })
        })
    })
    describe('I visit /admin/survey/new', function() {
        beforeEach(function() {
            cy.visit(`${Cypress.env('HOST')}/admin/survey/new/`, {
                auth: {
                    username: Cypress.env('ADMIN_USERNAME'),
                    password: Cypress.env('ADMIN_PASSWORD')
                }
            })
        })
        it('it displays the create survey form', function() {
            cy.get('.page-panel form').should('be.visible')
            cy.get('.page-panel form legend').should('contain', 'Create a New Survey')
            cy.get('.page-panel form').should('contain', 'Name')
            cy.get('.page-panel form').should('contain', 'Team')
            cy.get('.page-panel form button').should('be.visible').should('contain', 'Create Survey')
            cy.get('.page-panel form #survey-name').should('have.prop', 'hidden')
            cy.get('#team-new-wrapper').should('not.be.visible')
            cy.get('#team-name').should('contain', '» Create New Team')
        })

        describe('I select "» Create New Team"', function() {
            it('it displays the new form field', function() {
                cy.get('#team-name').select('» Create New Team')
                cy.get('#team-new-wrapper').should('be.visible')
                cy.get('#team-new').should('be.visible').should('have.prop', 'placeholder', 'New Team Name')
            })
        })

        describe('I deselect "» Create New Team"', function() {
            it('it displays the new form field', function() {
                cy.get('#team-name').select(Cypress.env('TEST_SURVEY_TEAM'))
                cy.get('#team-new-wrapper').should('not.be.visible')
                cy.get('#team-new').should('not.be.visible')
                    .should('have.prop', 'placeholder', 'New Team Name')
                    .should('have.value', '')
            })
        })

        it('it displays site navigation', function() {
            cy.get('.navbar').should('be.visible')
            cy.get('.navbar').contains("Dashboard").should('not.have.class', 'active')
            cy.get('.navbar').contains("Surveys").should('have.class', 'active')
        })
        it('it displays sub-navigation', function() {
            cy.get('.sub-navbar').should('be.visible')
            cy.get('.sub-navbar form').should('not.exist')
            cy.get('.sub-navbar').contains('All Surveys').should('not.have.class', 'active')
            cy.get('.sub-navbar').contains('Create New').should('have.class', 'active')
        })
        describe('I click on the All Surveys sub-navigation link', function() {
            it('it takes me to the the survey list', function() {
                cy.get('.sub-navbar').contains('All Surveys').click()
                cy.url().should('equal', `${Cypress.env('HOST')}/admin/survey/`)
            })
        })
        describe('I click on the Create New sub-navigation link', function() {
            it('it takes me to the the survey list', function() {
                cy.get('.sub-navbar').contains('Create New').click()
                cy.url().should('equal', `${Cypress.env('HOST')}/admin/survey/new/`)
            })
        })
        describe('I click on the dashboard navigation link', function() {
            it('it takes me to the team dashboard', function() {
                cy.get('.navbar').contains('Dashboard').click()
                cy.url().should('equal', `${Cypress.env('HOST')}/admin/dashboard/`)
            })
        })
        describe('I click on the surveys navigation link', function() {
            it('it takes me to the survey list', function() {
                cy.get('.navbar').contains('Surveys').click()
                cy.url().should('equal', `${Cypress.env('HOST')}/admin/survey/`)
            })
        })
    })
})