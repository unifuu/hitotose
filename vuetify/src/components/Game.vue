<template>
  <v-container>
    <v-card>
      <v-layout>
        <v-app-bar color="purple-darken-3" image="">
          <template v-slot:image>
            <v-img gradient="to top right, rgba(19,84,122,0.4), rgba(128,208,199,0.3)"></v-img>
          </template>

          <v-app-bar-nav-icon variant="text" @click.stop="drawer = !drawer"></v-app-bar-nav-icon>

          <v-app-bar-title>Library</v-app-bar-title>

          <v-spacer></v-spacer>

          <v-text-field class="mx-5" style="max-width: 300px" clearable prepend-inner-icon="mdi-magnify"
            density="comfortable" variant="solo" placeholder="Search..." hide-details v-model="keyword"
            @keyup.enter="performSearch" @click:clear="clearSearch" />

          <v-btn-toggle v-model="status" class="pr-3" rounded="xl" borderless divided>
            <v-btn value="Played">
                <v-icon :color="color('Played')" size="30">mdi-progress-check</v-icon>
            </v-btn>
            <v-btn value="Playing">
                <v-icon :color="color('Playing')" size="30">mdi-progress-star-four-points</v-icon>
            </v-btn>
            <v-btn value="ToPlay">
                <v-icon :color="color('ToPlay')" size="30">mdi-progress-close</v-icon>
            </v-btn>
          </v-btn-toggle>
        </v-app-bar>

        <v-main>
          <v-container fluid>
            <v-card>
              <v-layout>
                <v-navigation-drawer v-model="drawer" expand-on-hover rail>
                  <v-list>
                    <v-list-item v-if="token" link prepend-icon="mdi-database-plus-outline" title="Create Game"
                      @click="createGameDialog = true" />
                  </v-list>

                  <v-divider v-if="token"></v-divider>

                  <v-list density="compact" nav>
                    <v-list-item prepend-icon="mdi-select-all" title="All" @click="platform = 'All'"
                      :active="platform === 'All'" />
                    <v-list-item prepend-icon="mdi-alien" title="PC" @click="platform = 'PC'"
                      :active="platform === 'PC'" />
                    <v-list-item prepend-icon="mdi-sony-playstation" title="Playstation"
                      @click="platform = 'PlayStation'" :active="platform === 'PlayStation'" />
                    <v-list-item prepend-icon="mdi-nintendo-switch" title="Nintendo Switch"
                      @click="platform = 'Nintendo Switch'" :active="platform === 'Nintendo Switch'" />
                    <v-list-item prepend-icon="mdi-microsoft-xbox" title="Xbox" @click="platform = 'Xbox'"
                      :active="platform === 'Xbox'" />
                    <v-list-item prepend-icon="mdi-tablet-cellphone" title="Mobile" @click="platform = 'Mobile'"
                      :active="platform === 'Mobile'" />
                  </v-list>
                </v-navigation-drawer>

                <v-main class="ma-1" style="min-height: 350px">
                  <div class="d-flex flex-wrap ga-5">
                    <v-card v-for="g in filteredGames(status)" :key="g.id" max-height="375" max-width="250"
                      variant="tonal" class="text-center">
                      <v-card-title style="word-break: wrap-all; white-space: normal; overflow: hidden; font-size: 15px; line-height: 1.5em; height: 3.5em;">{{ g.title }}</v-card-title>
                      <v-card-subtitle>{{ hourOfPlayedTime(g.played_time) }}h {{ minOfPlayedTime(g.played_time)
                        }}m</v-card-subtitle>
                      <v-img cover :aspect-ratio="1" width="250"
                        :src="`http://localhost:8080/assets/images/games/${g.id}.webp`"></v-img>
                      <v-card-actions>
                        <v-btn variant="outlined" density="comfortable" :icon="icon(g.status)" :color="color(g.status)"
                          @click="fetchUpdatingGame(g.id)" />
                        <v-btn variant="outlined" density="comfortable" :icon="icon(g.platform)" :color="color(g.platform)" />
                        <v-btn variant="outlined" density="comfortable" icon @click="openUpdateGameRatingDialog(g.id, g.title, g.rating)">
                          <v-progress-circular :model-value="percentage(g.rating, 10)">
                            <template v-slot:default>
                              <span style="font-size: 0.7rem;">{{ g.rating ? g.rating : '-' }}</span>
                            </template>
                          </v-progress-circular>
                        </v-btn>
                      </v-card-actions>
                    </v-card>
                  </div>
                </v-main>
              </v-layout>
            </v-card>
          </v-container>

          <v-pagination v-model="currentPage" :length="totalPage" :total-visible="4" />
        </v-main>
      </v-layout>
    </v-card>

    <!-- Create Game Dialog -->
    <div class="pa-4 text-center">
      <v-dialog v-model="createGameDialog" max-width="600">
        <v-card prepend-icon="mdi-database-plus-outline" title="Create Game">
          <v-form method="post" action="/api/game/create">
            <v-card-text>
              <v-row dense>
                <v-col cols="12">
                  <v-text-field name="title" label="Title" required />
                </v-col>

                <v-col cols="12" md="6" sm="6">
                  <v-text-field name="developer" label="Developer" />
                </v-col>

                <v-col cols="12" md="6" sm="6">
                  <v-text-field name="publisher" label="Publisher" />
                </v-col>

                <v-col cols="12" sm="6">
                  <v-select name="genre" v-model="selectedGenre" :items="genres" item-title="label" item-value="value"
                    label="Genre" required />
                </v-col>

                <v-col cols="12" sm="6">
                  <v-select name="platform" v-model="selectedPlatform" :items="platforms" item-title="label"
                    item-value="value" label="Platform" required />
                </v-col>
              </v-row>
            </v-card-text>

            <v-divider></v-divider>

            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn type="button" text="Close" variant="plain" @click="createGameDialog = false" />
              <v-btn type="submit" color="success" text="Save" variant="tonal" />
            </v-card-actions>
          </v-form>
        </v-card>
      </v-dialog>
    </div>

    <!-- Update Game Dialog -->
    <div class="pa-4 text-center">
      <v-dialog v-model="updateGameDialog" max-height="600" max-width="600">
        <v-form method="post" encType="multipart/form-data" action="/api/game/update">
          <v-card prepend-icon="mdi-database-edit-outline" title="Edit Game">
            <v-card-text>
              <v-row dense>
                <v-col cols="12">
                  <v-text-field variant="solo-filled" name="id" label="ID" v-model="updatingGame.id" readonly />
                </v-col>

                <v-col cols="12" class="dense">
                  <v-text-field variant="solo-filled" name="title" label="Title" v-model="updatingGame.title"
                    required />
                </v-col>

                <v-col cols="12" sm="6" class="dense">
                  <v-text-field variant="solo-filled" name="developer" label="Developer"
                    v-model="updatingGame.developer" />
                </v-col>

                <v-col cols="12" sm="6" class="dense">
                  <v-text-field variant="solo-filled" name="publisher" label="Publisher"
                    v-model="updatingGame.publisher" />                  
                </v-col>

                <v-col cols="12" sm="6" class="dense">
                  <v-select variant="solo-filled" name="genre" v-model="updatingGame.genre" :items="genres"
                    item-title="label" item-value="value" label="Genre" required />
                </v-col>

                <v-col cols="12" sm="6" class="dense">
                  <v-select variant="solo-filled" name="platform" v-model="updatingGame.platform" :items="platforms"
                    item-title="label" item-value="value" label="Platform" required />
                </v-col>

                <v-col cols="12" sm="6" class="dense">
                  <v-select variant="solo-filled" name="status" v-model="updatingGame.status" :items="statuses"
                    item-title="label" item-value="value" label="Status" required />
                </v-col>

                <v-col cols="12" sm="6" class="dense">
                  <v-select variant="solo-filled" name="rating" v-model="updatingGame.rating" :items="ratings"
                    item-title="label" item-value="value" label="Rating" required />
                </v-col>

                <v-col class="dense">
                  <v-file-input v-model="coverFile" name="cover" label="Cover" variant="solo-filled"></v-file-input>
                  <!-- <input type="file" name="cover" /> -->
                </v-col>
              </v-row>
            </v-card-text>

            <v-divider></v-divider>

            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn type="button" text="Close" variant="plain" @click="updateGameDialog = false"></v-btn>
              <v-btn type="submit" color="success" text="Save" variant="tonal"
                @click="updateGameDialog = false"></v-btn>
            </v-card-actions>
          </v-card>
        </v-form>
      </v-dialog>

    </div>

    <!-- Update Rating Dialog -->
    <div class="pa-4 text-center">
      <v-dialog v-model="updateGameRatingDialog" max-height="600" max-width="450">
        <v-form method="post" encType="multipart/form-data" action="/api/game/update/rating">
          <v-card prepend-icon="mdi-database-edit-outline" title="Edit Rating">
            <v-card-text>
              <v-row dense>
                <v-col cols="12">
                  <v-text-field variant="solo-filled" name="id" label="ID" v-model="updatingGame.id" readonly />
                </v-col>

                <v-col cols="12" class="dense">
                  <v-text-field variant="solo-filled" name="title" label="Title" v-model="updatingGame.title" readonly />
                </v-col>

                <v-col cols="12" class="dense">
                  <v-select variant="solo-filled" name="rating" v-model="updatingGame.rating" :items="ratings"
                    item-title="label" item-value="value" label="Rating" required />
                </v-col>
              </v-row>
            </v-card-text>

            <v-divider></v-divider>

            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn type="button" text="Close" variant="plain" @click="updateGameDialog = false"></v-btn>
              <v-btn type="submit" color="success" text="Save" variant="tonal" @click="updateGameRatingDialog = false"></v-btn>
            </v-card-actions>
          </v-card>
        </v-form>
      </v-dialog>
    </div>

  </v-container>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import optionJson from '../assets/selectOptions.json'
import { getCookie } from './token'

const selectOptions: SelectOptions = optionJson as unknown as SelectOptions
const loading = ref<boolean>(false)
const status = ref<string>('Playing')
const platform = ref<string>('All')
const drawer = ref<boolean>(false)
const createGameDialog = ref<boolean>(false)
const updateGameDialog = ref<boolean>(false)
const updateGameRatingDialog = ref<boolean>(false)
const updatingGame = ref<Game>(newGame())
const games = ref<Game[]>([])
const error = ref<string | null>(null)
const token = ref<string | null>(null)
const keyword = ref<string | null>()
const currentPage = ref<number>(1)
const totalPage = ref<number>(1)
const coverFile = ref<File | null>(null)

// Select
const selectedGenre = ref<string>('')
const selectedPlatform = ref<string>('')
const genres = ref<Option[]>(selectOptions.genres)
const platforms = ref<Option[]>(selectOptions.platforms)
const ratings = ref<Option[]>(selectOptions.ratings)
const statuses = ['Played', 'Playing', 'ToPlay']

interface Game {
  id: string;
  title: string;
  genre: string;
  platform: string;
  developer: string;
  publisher: string;
  status: string;
  rating: number;
  played_time: number;
}

function newGame(): Game {
  return {
    id: '',
    title: '',
    genre: '',
    platform: '',
    developer: '',
    publisher: '',
    status: '',
    rating: 0,
    played_time: 0,
  }
}

export interface Option {
  value: string
  label: string
}

interface SelectOptions {
  genres: Option[]
  platforms: Option[]
  ratings: Option[]
}

watch([status, platform, currentPage], () => {
  loading.value = true
  fetchGames()
})

onMounted(() => {
  checkToken()
  fetchGames()
})

const icon = (key: string) => {
  switch (key) {
    // Statuses
    case 'Played':
      return 'mdi-progress-check'
    case 'Playing':
      return 'mdi-progress-star-four-points'
    case 'ToPlay':
      return 'mdi-progress-close'

    // Platforms
    case 'PC':
      return 'mdi-alien'
    case 'PlayStation':
      return 'mdi-sony-playstation'
    case 'Nintendo Switch':
      return 'mdi-nintendo-switch'
    case 'Xbox':
      return 'mdi-microsoft-xbox'
    case 'Mobile':
      return 'mdi-tablet-cellphone'

    default:
      return ''
  }
}

const color = (key: string) => {
  switch (key) {
    // Statuses
    case 'Played':
      return '#FFFFFF'
    case 'Playing':
      return '#81C784'
    case 'ToPlay':
      return '#FFF59D'

    // Platforms
    case 'PC':
      return '#FFFE3F'
    case 'PlayStation':
      return '#2E6DB4'
    case 'Nintendo Switch':
      return '#E60012'
    case 'Xbox':
      return '#107C10'
    case 'Mobile':
      return '#A6A8AB'
  }
}

const percentage = (current: number, total: number) => {
  switch (total) {
    case 0:
      return 0
    default:
      if (current >= total) {
        return 100
      } else {
        return Math.round((current / total) * 100)
      }
  }
}

const openUpdateGameRatingDialog = (id: string, title: string, rating: number) => {
  updatingGame.value.id = id
  updatingGame.value.title = title
  updatingGame.value.rating = rating
  updateGameRatingDialog.value = true
}

const clearSearch = () => {
  keyword.value = ""
  performSearch()
}

const performSearch = async () => {
  status.value = ""
  selectedPlatform.value = "All"
  fetchGames()
}

const fetchUpdatingGame = async (id: string) => {
  try {
    const response = await fetch(`/api/game/update?id=${id}`)

    if (!response.ok) {
      throw new Error('Network response was not ok')
    }

    const responseBody = await response.text()
    if (responseBody) {
      const data = JSON.parse(responseBody)

      if (data && data.game) {
        updatingGame.value = data.game
      }

      updateGameDialog.value = true
    } else {
      throw new Error('Empty response body')
    }
  } catch (err: any) {
    error.value = err.message
    console.log(error.value)
  } finally {
    loading.value = false
  }
}

const filteredGames = (status: string | undefined) => {
  if (!games?.value) { return [] }
  if (status === undefined || status === "") { return games.value }
  return games.value.filter(game => game.status === status)
}

const fetchGames = async () => {
  try {
    const baseURL = `/api/game/pages?page=${currentPage.value}`
    let queryParam = status.value ? `&status=${status.value}` : ''
    queryParam += keyword.value ? `&keyword=${keyword.value}` : ''
    const response = await fetch(`${baseURL}${queryParam}`)
    if (!response.ok) {
      throw new Error('Network response was not ok')
    }
    const data = await response.json()
    games.value = data.games as Game[]
    totalPage.value = data.total_page
  } catch (err: any) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

function hourOfPlayedTime(playedTime: number): number {
  const hour = Math.floor(playedTime / 60)
  return hour
}

function minOfPlayedTime(playedTime: number): number {
  const min = playedTime % 60
  return min
}

const checkToken = () => {
  // token.value = localStorage.getItem('auth_token')
  token.value = getCookie("auth_token")
}
</script>

<style scoped>
.dense {
  margin-top: -1rem
}
</style>