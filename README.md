# Gobus
A bus you didn't know you needed with Go

Allows for creating an event like architecture with upfront configuration and clean decoupling of code through events.

Event functions can be run asynchronously or synchronously and the calling function will receive all of the errors that
happened in order to handle the errors properly.



# Simple Usage

1. Define your events

    ```go
    var (
        USER_SIGNED_UP = gobus.Event("user:signed-up")
        USER_LOGGED_OUT = gobus.Event("user:logged-out")
    )
    ```
2. Create your event functions -- these functions should be glue code between your systems

    You wan't your internal logic to deal with things that make sense to your application.

    ```
    Things like:
        - user id
        - email
        - website url
    ```
    Parse the args in this function and call the business logic function with

    ```go
        import "your/internal/email"
        import "your/internal/cache"

        // Your own user definition OR a glue definition
        //
        // Really up to you.
        type User struct {
            ID    int64
            Email string

        }

        // Sends a welcome email
        func EventUserSignedUpSendWelcomeEmail(args any) error {
            if _, ok := args.(*User); ok {
                email.SendUserWelcomeEmail(args.(*User))
            } else {
                return NewInvalidArgError{"EventUserSignedUpSendWelcomeEmail", User{}, args}
            }
        }

        // Sends a verification email
        func EventUserSignedUpSendVerificationEmail(args *any) error {
            if _, ok := args.(*User); ok {
                email.SendUserWelcomeEmail(args.(*User))
            } else {
                return gobus.NewInvalidArgError{"EventUserSignedUpSendVerificationEmail", User{}, args}
            }
        }

        // Clears a user cache
        func EventUserLoggedOutClearUserCache(args *args) error {
            if _, ok := args.(*User); ok {
                cache.ClearUser(args.(*User))
            } else {
                return gobus.NewInvalidArgError{"EventUserLoggedOutClearUserCache", User{}, args}
            }
        }
    ```
3. Register your event functions

    ```go
        gobus.SetEventFuncs(gobus.EventFuncsMap{
            USER_SIGNED_UP: []gobus.EventFunc{
                EventUserSignedUpSendWelcomeEmail,
                EventUserSignedUpSendVerificationEmail,
            },
            USER_LOGGED_OUT: []EventFunc {
                EventUserLoggedOutClearUserCache,
            },
	    })
    ```