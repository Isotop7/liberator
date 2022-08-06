<script>
    import Navbar from "./Navbar.svelte";
</script>

<Navbar />

<!-- This example requires Tailwind CSS v2.0+ -->
<div class="py-12 bg-white">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="lg:text-center">
            <h2
                class="text-base text-indigo-600 font-semibold tracking-wide uppercase"
            >
                liberator
            </h2>
            <p
                class="mt-2 text-3xl leading-8 font-extrabold tracking-tight text-gray-900 sm:text-4xl"
            >
                Luk @ ma bukz, bruv.
            </p>
            <p class="mt-4 max-w-2xl text-2xl text-gray-500 lg:mx-auto">
                <span style="italic">liberator</span> is a dead simple web app for
                your digital book shelve. Helps you to keep track of them books.
            </p>
            <p class="mt-4 max-w-xl text-xl text-gray-500 lg:mx-auto">
                It's written in Go and uses <a
                    href="https://github.com/livebud/bud">Bud</a
                >.
            </p>
        </div>
    </div>
</div>
