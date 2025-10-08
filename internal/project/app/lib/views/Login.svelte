<script lang="ts">
    import Layout from "$lib/components/Layout.svelte"
    import SuperForm from "$lib/components/forms/SuperForm.svelte"
    import Field from "$lib/components/forms/Field.svelte"
    import Control from "$lib/components/forms/Control.svelte"
    import Label from "$lib/components/forms/Label.svelte"
    import FieldErrors from "$lib/components/forms/FieldErrors.svelte"
    import Description from "$lib/components/forms/Description.svelte"
    import { superForm, required, email, minLength } from "$lib/scripts/core/form.svelte.ts"

    // Props type will be generated from Go
    // type Props = { form: FormState<{email: string, password: string}> }
    let { form: serverForm }: any = $props()

    // Initialize form with validation rules
    const form = superForm(
        {
            email: serverForm?.data?.email || "",
            password: serverForm?.data?.password || "",
        },
        {
            email: [required(), email()],
            password: [required(), minLength(8)],
        }
    )

    // Update form state from server if there are errors
    if (serverForm?.errors) {
        form.updateFromServer(serverForm)
    }
</script>

<Layout title="Login">
    <div class="w-full max-w-md mx-auto">
        <div class="card bg-base-200 shadow-xl">
            <div class="card-body">
                <h2 class="card-title text-2xl mb-4">Login to Your Account</h2>

                <SuperForm {form} action="/login" method="POST">
                    {#snippet children({ pending, valid, message })}
                        {#if message}
                            <div class="alert alert-error mb-4">
                                <span>{message}</span>
                            </div>
                        {/if}

                        <Field {form} name="email">
                            {#snippet children({ hasError })}
                                <div class="form-control w-full">
                                    <Label class="label">
                                        <span class="label-text">Email</span>
                                    </Label>
                                    <Control>
                                        {#snippet children({ props })}
                                            <input
                                                {...props}
                                                type="email"
                                                bind:value={form.state.data.email}
                                                class="input input-bordered w-full"
                                                class:input-error={hasError}
                                                placeholder="you@example.com"
                                            />
                                        {/snippet}
                                    </Control>
                                    <Description>
                                        Enter your email address
                                    </Description>
                                    <FieldErrors />
                                </div>
                            {/snippet}
                        </Field>

                        <Field {form} name="password" class="mt-4">
                            {#snippet children({ hasError })}
                                <div class="form-control w-full">
                                    <Label class="label">
                                        <span class="label-text">Password</span>
                                    </Label>
                                    <Control>
                                        {#snippet children({ props })}
                                            <input
                                                {...props}
                                                type="password"
                                                bind:value={form.state.data.password}
                                                class="input input-bordered w-full"
                                                class:input-error={hasError}
                                                placeholder="Enter your password"
                                            />
                                        {/snippet}
                                    </Control>
                                    <Description>
                                        Must be at least 8 characters
                                    </Description>
                                    <FieldErrors />
                                </div>
                            {/snippet}
                        </Field>

                        <div class="form-control mt-6">
                            <button
                                type="submit"
                                class="btn btn-primary"
                                disabled={pending || !valid}
                            >
                                {#if pending}
                                    <span class="loading loading-spinner"></span>
                                    Logging in...
                                {:else}
                                    Login
                                {/if}
                            </button>
                        </div>

                        <div class="divider">OR</div>

                        <a href="/welcome" class="btn btn-ghost btn-sm">
                            Back to Welcome
                        </a>
                    {/snippet}
                </SuperForm>
            </div>
        </div>
    </div>
</Layout>
