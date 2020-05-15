<template>
  <div class="home page">
    <h1>{{ title }}</h1>
   
    <UiModal closeOnOverlay :show.sync="isShownModal">
      <div class="some-modal-content">
        hi here
        <div class="buttons">
          <button @click="submitModalHandler">ok</button>
        </div>
      </div>
    </UiModal>

    <div class="form">
      Routes
      <select ref="route" v-model="selectedRoute" @change="showRoute">
        <option v-for="route in routes" :key="route" :id="route">{{ route }}</option>
      </select>
    </div>

    <button class="button" @click="startRoute">Start</button>

    <div class="map-container"
      :style="{minHeight: height + 'px'}"
    >
        <div ref="map"
            class="map"
        />
        <div class="map-hidden">
            <slot />
        </div>
        <slot name="visible" />
    </div>
  </div>
</template>

<script>
import UiModal from '@/components/UiModal.vue'
import MapLoader from '@/lib/map-loader.js'
import axios from 'axios'

export default {
  name: 'IndexPage',

  components: {
    UiModal
  },

  data () {
    return {
      msg: 'Route Simulator',
      isShownModal: false,
      inputError: false,
      checkboxState: false,

      // google map keys
      apiKey: 'AIzaSyAmeDzihSpzZ0QGae4UUoNLKRQPl7NpvVg',
      version: 'quarterly',
      libraries: ['geometry'],

      center: {
        lat: -7.931656836523958,
        lng: 112.61701910470126
      },
      zoom: 14,
      mapType: 'roadmap',

      routes: [],
      selectedRoute: '',
      
      traveler1: [],
      traveler2: []
    }
  },

  computed: {
    height () {
      return window.innerHeight - 250
    }
  },

  watch: {
    height () {
      if (!this._map) return

      window.google.maps.event.trigger(this._map, 'resize')
    }
  },

  mounted () {
    MapLoader.load(this.apiKey, this.version, this.libraries)
    MapLoader.ensureReady().then(() => {
      this.initMap()
      this.loadRoutes()
    })
  },

  destroyed () {
    if (this._map) {
      google.maps.event.clearListeners(this._map, 'bounds_changed')
      google.maps.event.clearListeners(this._map, 'zoom_changed')
      google.maps.event.clearListeners(this._map, 'click')
      google.maps.event.clearListeners(this._map, 'idle')
      this._map = null
    }
  },

  methods: {

    initMap () {
      if (this._map) {
        this.$emit('ready', this._map)
        return
      }

      setTimeout(() => {
        let options = Object.assign({
          center: this.center,
          zoom: this.zoom,
          mapTypeId: this.mapType,
          mapTypeControlOptions: {
            style: google.maps.MapTypeControlStyle.HORIZONTAL_BAR,
            position: google.maps.ControlPosition.TOP_RIGHT
          },
          zoomControl: true,
          zoomControlOptions: {
            position: google.maps.ControlPosition.TOP_LEFT
          },
          scaleControl: true,
          streetViewControl: true,
          streetViewControlOptions: {
            position: google.maps.ControlPosition.TOP_LEFT
          },
          disableDoubleClickZoom: true,
          panControl: true,
          styles: [
            {
              featureType: 'poi',
              stylers: [{ visibility: 'off' }]
            }
          ]
        }, this.options || {})

        console.log('ref', this.$refs)
        this._map = new google.maps.Map(this.$refs.map, options)
        google.maps.event.trigger(this._map, 'resize')

        this.routeMap = new google.maps.Polyline({
          strokeColor: '#00FF00',
          strokeOpacity: 1.0,
          strokeWeight: 3,
          geodesic: true,
          map: this._map
        })

        this.traveler1Map = new google.maps.Polyline({
          strokeColor: '#FF0000',
          strokeOpacity: 1.0,
          strokeWeight: 3,
          geodesic: true,
          map: this._map
        })

        this.traveler2Map = new google.maps.Polyline({
          strokeColor: '#0000FF',
          strokeOpacity: 1.0,
          strokeWeight: 3,
          geodesic: true,
          map: this._map
        })

        google.maps.event.addListener(
          this._map,
          'bounds_changed',
          e => {
            this.$emit('boundsChanged', this._map.getBounds())
          }
        )

        google.maps.event.addListener(this._map, 'zoom_changed', e => {
          this.$emit('zoomChanged', this._map.getZoom())
        })

        google.maps.event.addListener(this._map, 'idle', e => {
          this.$emit('idle', this._map)
        })

        google.maps.event.addListener(this._map, 'click', e => {
          const pos = {
            lat: e.latLng.lat(),
            lng: e.latLng.lng()
          }
                              
          this.$emit('mapClicked', pos)
        })

        // notify parent component that map is ready
        this.$emit('ready', this._map)

      }, 500)
    },

    loadRoutes () {
      this.busy = true
      axios.get('/api/routes')
        .then(resp => {
          this.routes = resp.data || []
          this.busy = false
        })
        .catch(e => {
          console.log(e)
          this.routes = []
          this.busy = false
        })
    },

    showRoute () {
      this.busy = true

      axios.get(`/api/routes/${this.selectedRoute}/coordinates`)
        .then(resp => {
          this.routeMap.setPath(resp.data || [])
          this.busy = false
        })
        .catch(e => {
          console.log(e)
          this.routeMap.setPath([])
          this.busy = false
        })
    },

    async getNextRoute (currentPos, reverse) {
      const params = {
        direction: reverse ? 'reverse' : '',
        current: JSON.stringify(currentPos)
      }

      return axios.get(`/api/routes/${this.selectedRoute}/next`, { params })
        .then(resp => {
          return resp.data
        })
    },

    async getStartRoute (reverse) {
      const params = {
        direction: reverse ? 'reverse' : ''
      }

      return axios.get(`/api/routes/${this.selectedRoute}/start`, { params })
        .then(resp => {
          return resp.data
        })
    },

    async traveler1Move () {
      let nextPost = await this.getNextRoute(this.traveler1Pos)
      if (!nextPost) return

      this.traveler1.push(this.traveler1Pos)
      this.traveler1Pos = nextPost

      this.traveler1Map.setPath(this.traveler1)

      let vm = this
      setTimeout(() => {
        vm.traveler1Move()
      }, 1000)
    },

    async traveler2Move () {
      let nextPost = await this.getNextRoute(this.traveler2Pos, true)
      if (!nextPost) return

      this.traveler2.push(this.traveler2Pos)
      this.traveler2Pos = nextPost

      this.traveler2Map.setPath(this.traveler2)

      let vm = this
      setTimeout(() => {
        vm.traveler2Move()
      }, 1000)
    },

    async startRoute () {
      this.traveler1 = []
      this.traveler2 = []
      this.traveler1Map.setPath([])
      this.traveler2Map.setPath([])

      this.traveler1Pos = await this.getStartRoute()
      this.traveler2Pos = await this.getStartRoute(true)

      this.traveler1Move()
      this.traveler2Move()
    },

    isMeetingPoint() {
      // TODO: calculate whether current traveler1 and traveler 2 position is meeting point, show information to UI
    },

    showToast () {
      this.$store.commit('toast/NEW', { type: 'success', message: 'hello' })
    },

    submitModalHandler () {
      // some logic
      this.isShownModal = false
    }
  }
}
</script>

<style lang="scss" scoped>
.some-modal-content {
  min-width: 400px;
  padding: 25px;

  .buttons button {
    padding: 10px;
    margin: 10px;
  }
}

.map-container {
    position: relative;
}

.map-container .map {
    left: 0;
    right: 0;
    top: 0;
    bottom: 0;
    position: absolute;
}
.map-hidden {
    display: none;
}

</style>
