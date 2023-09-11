# Gobus
A library that makes it easy to execute a distributed sequence of functions when an event gets emitted.

Simplifies creating event based architecture.

You can focus on:
- The input of event handlers
- Emitting events
- Handling errors


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
                return gobus.NewInvalidArgError{"EventUserSignedUpSendWelcomeEmail", User{}, args}
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


## Examples
There's a folder with several examples which can show how this library can be used to create distributed workflows.

