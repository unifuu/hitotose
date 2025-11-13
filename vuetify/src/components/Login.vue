<template>
    <v-container fluid fill-height class="login-container">
        <v-row align="center" justify="center">
            <v-col cols="12" sm="8" md="4">
                <v-card class="pa-4" outlined>
                    <v-card-text>
                        <v-form v-model="valid" @submit.prevent="handleLogin">
                            <v-text-field v-model="username" label="Username" required></v-text-field>
                            <v-text-field v-model="password" label="Password" type="password" required></v-text-field>
                            <v-btn outlined block type="submit" color="primary" height="45px">LOGIN</v-btn>
                        </v-form>
                    </v-card-text>
                </v-card>
            </v-col>
        </v-row>
    </v-container>
</template>

<script lang="ts" setup>
import router from '@/router';
import { ref } from 'vue'
import { VTextField, VBtn, VCard, VCardText, VContainer, VRow, VCol, VForm } from 'vuetify/components'

const username = ref('')
const password = ref('')
const token = ref<string | null>(null)
const valid = ref(false)

const handleLogin = async () => {
    try {
    const response = await fetch('http://127.0.0.1:8080/api/user/checkAuth', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username: username.value,
        password: password.value,
      }),
    });

    if (response.ok) {
      const data = await response.json();
      token.value = data.auth_token;

      if (token.value) {
        console.log(token.value)
        // localStorage.setItem('auth_token', token.value)
        setCookie("auth_token", token.value)
        router.push('/game')
      }
    } else {
      console.error('Login failed');
    }
  } catch (error) {
    console.error('An error occurred:', error);
  }
}

function setCookie(name: string, val: string) {
    const expire = new Date()
    expire.setTime(expire.getTime() + (1000 * 60 * 60 * 24 * 30)) // A month
    document.cookie = name+"="+val+"; expires="+expire.toUTCString()+"; path=/"
}
</script>

<style scoped>
.login-container {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100vh;
}

.v-card {
    max-width: 400px;
    width: 100%;
}
</style>