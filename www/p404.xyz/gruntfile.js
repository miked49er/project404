module.exports = function(grunt) {

  grunt.initConfig({
    jshint: {
      files: ["*.js", "lib/js/*.js"],
      options: {
        esnext: true,
        globals: {
          jQuery: true
        } // globals
      } // options
    }, // jshint
    uglify: {
      options: {
        mangle: {
          except: [
            '$timeout',
            '$eval',
            '$window',
            '$document'
          ]
        } // mangle
      }, // options
      my_target: {
        files: {
          'public/javascripts/script.min.js': ['lib/js/*.js']
        } // files
      } // my_target
    }, // uglify
    compass: {
      dev: {
        options: {
          config: 'config.rb'
        } // options
      } // dev
    }, // compass
    watch: {
      options: { livereload: true },
      scripts: {
        files: ['lib/js/*.js', '*.js'],
        tasks: ['uglify', 'jshint']
      }, // scripts
      sass: {
        files: ['lib/sass/*.scss'],
        tasks: ['compass:dev']
      }, // sass
      html: {
        files: ['templates/*.html']
      } // html
    } // watch
  }); // initConfig

  grunt.loadNpmTasks('grunt-contrib-uglify');
  grunt.loadNpmTasks('grunt-contrib-watch');
  grunt.loadNpmTasks('grunt-contrib-compass');
  grunt.loadNpmTasks("grunt-contrib-jshint");

  grunt.registerTask('js', ['uglify']);
  grunt.registerTask('css', ['compass:dev']);
  grunt.registerTask('default', ['watch']);
};
