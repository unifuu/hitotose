<template>
  <div>
    <h1>Titles</h1>
    <table>
      <thead>
        <tr>
          <th>Title</th>
          <th>Genre</th>
          <th>Platform</th>
          <th>Developer</th>
          <th>Publisher</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="game in games" :key="game.id">
          <td>{{ game.title }}</td>
          <td>{{ game.genre }}</td>
          <td>{{ game.platform }}</td>
          <td>{{ game.developer }}</td>
          <td>{{ game.publisher }}</td>
        </tr>
      </tbody>
    </table>
    <div v-if="loading">Loading...</div>
    <div v-if="error">{{ error }}</div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      games: [],
      loading: true,
      error: null,
    };
  },
  mounted() {
    this.fetchGames();
  },
  methods: {
    async fetchGames() {
      try {
        const response = await fetch('http://127.0.0.1:8080/api/game/pages');
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        const data = await response.json();
        this.games = data.games; // Access the 'games' array from the response
      } catch (error) {
        this.error = error.message;
      } finally {
        this.loading = false;
      }
    },
  },
};
</script>

<style>
table {
  width: 100%;
  border-collapse: collapse;
}
th, td {
  border: 1px solid #ddd;
  padding: 8px;
}
th {
  background-color: #f2f2f2;
}
</style>
