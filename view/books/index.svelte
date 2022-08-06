<script>
    import Navbar from "../Navbar.svelte";
    export let books = [];
</script>

<Navbar />

<h1>List of books</h1>

<ul>
    {#each books as book}
        <li>ID: {books.id}</li>
        <li>Title: {book.title}</li>
        <li>Author: {book.author}</li>
        <li>Language: {book.language}</li>
        <li>Category: {book.category}</li>
        <li>ISBN10: {book.isbn10}</li>
        <li>ISBN13: {book.isbn13}</li>
        <li>PageCount: {book.pagecount}</li>
        <li>Rating: {book.rating}</li>
    {/each}
</ul>
