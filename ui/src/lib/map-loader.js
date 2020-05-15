// taken from github.com/Akryum/vue-googlemaps/blob/master/src/lib-loader.js
const loader = {
    loaded: false,
    readyPromises: [],

    /**
     * @param apiKey    API Key, or object with the URL parameters. For example
     *                  to use Google Maps Premium API, pass
     *                    `{ client: <YOUR-CLIENT-ID> }`.
     *                  You may pass the libraries and/or version (as `v`) parameter into
     *                  this parameter and skip the next two parameters
     * @param version   Google for Maps version
     * @param libraries Libraries to load (@see
     *                  https://developers.google.com/maps/documentation/javascript/libraries)
     */
    load(apiKey, version, libraries) {
        // prevent burst call
        if (this._tid) {
            clearTimeout(this._tid)
        }

        this._tid = setTimeout(() => {
            this._load(apiKey, version, libraries)
        }, 500)
    },

    _load(apiKey, version, libraries) {
        if (typeof window === 'undefined') {
            // Do nothing if run from server-side
            return Promise.resolve()
        }
        if (!this.loaded && (!window.google || !window.google.maps)) {
            const script = document.createElement('script')
            let options = {}

            if (typeof apiKey === 'string') {
                options.key = apiKey
            } else if (typeof apiKey === 'object') {
                options = Object.assign(options, apiKey)
            } else {
                Promise.reject(new Error('apiKey should either string or object'))
            }

            if (libraries && libraries.length) {
                options.libraries = libraries.join(',')
            }

            if (version) {
                options.v = version
            }

            const url = 'https://maps.googleapis.com/maps/api/js?' +
                Object.keys(options)
                    .map(k => encodeURIComponent(k) + '=' + encodeURIComponent(options[k]))
                    .join('&')

            script.src = url
            script.type = 'application/javascript'

            script.onreadystatechange = () => {
                if (script.readyState === 'complete') {
                    resolve()
                }
            }

            script.onload = () => {
                this._setLoaded()
            }

            const head = document.getElementsByTagName('head')[0]
            head.appendChild(script)
        } else {
            console.warn('The Google Maps library is already loaded')
            this._setLoaded()
        }
    },

    ensureReady() {
        if (this.loaded) {
            return Promise.resolve()
        } else {
            const promise = new Promise((resolve) => {
                this.readyPromises.push(resolve)
            })
            return promise
        }
    },

    _setLoaded() {
        this.loaded = true
        for (const resolve of this.readyPromises) {
            resolve()
        }
        this.readyPromises = []
    }
}

export default loader

