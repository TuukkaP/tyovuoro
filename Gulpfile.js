/* ----------------------------------------------------------------------------
 * Imports
 * ------------------------------------------------------------------------- */
 
var gulp       = require("gulp");
var addsrc     = require("gulp-add-src");
var args       = require("yargs").argv;
var autoprefix = require("autoprefixer-core");
var cache      = require("gulp-cached");
var clean      = require("del");
var collect    = require("gulp-rev-collector");
var concat     = require("gulp-concat");
var ignore     = require("gulp-ignore");
var mincss     = require("gulp-minify-css");
var minhtml    = require("gulp-htmlmin");
var modernizr  = require("gulp-modernizr");
var mqpacker   = require("css-mqpacker");
var notifier   = require("node-notifier");
var gulpif     = require("gulp-if");
var pixrem     = require("pixrem");
var plumber    = require("gulp-plumber");
var postcss    = require("gulp-postcss");
var reload     = require("gulp-livereload");
var rev        = require("gulp-rev");
var sass       = require("gulp-sass");
var sourcemaps = require("gulp-sourcemaps");
var sync       = require("gulp-sync")(gulp).sync;
var child      = require("child_process");
var uglify     = require("gulp-uglify");
var util       = require("gulp-util");
var vinyl      = require("vinyl-paths");
var debug      = require("gulp-debug");
var react      = require('gulp-react');
var order      = require('gulp-order');
 
/* ----------------------------------------------------------------------------
 * Locals
 * ------------------------------------------------------------------------- */
 
/* Application server */
var server = null;
 
/* ----------------------------------------------------------------------------
 * Overrides
 * ------------------------------------------------------------------------- */
 
/*
 * Override gulp.src() for nicer error handling.
 */
var src = gulp.src;
gulp.src = function() {
  return src.apply(gulp, arguments)
    .pipe(plumber(function(error) {
      util.log(util.colors.red(
        "Error (" + error.plugin + "): " + error.message
      ));
      notifier.notify({
        title: "Error (" + error.plugin + ")",
        message: error.message.split("\n")[0]
      });
      this.emit("end");
    })
  );
};
 
/* ----------------------------------------------------------------------------
 * Assets pipeline
 * ------------------------------------------------------------------------- */
 
/*
 * Build stylesheets from SASS source.
 */
gulp.task("app:stylesheets", function() {
  return gulp.src("app/styles/*.scss")
    .pipe(sass({
      includePaths: [
        "bower_components/foundation/scss"
      ],
      outputStyle: 'nested',
      errLogToConsole: true}))
    .pipe(gulpif(args.sourcemaps, sourcemaps.init()))
    .pipe(gulpif(args.production,
      postcss([
        autoprefix(),
        mqpacker,
        pixrem("10px")
      ])))
    .pipe(gulpif(args.sourcemaps, sourcemaps.write()))
    .pipe(gulpif(args.production, mincss()))
    .pipe(gulp.dest("public/static/"))
    .pipe(reload());
});
 
/*
 * Build javascripts from Bower components and source.
 */
gulp.task("app:javascripts", function() {
  return gulp.src([
    "bower_components/jquery/dist/jquery.js",
    "bower_components/react/JSXTransformer.js",
    "bower_components/react/react-with-addons.js",
    "bower_components/react-router/build/umd/ReactRouter.js",
    "bower_components/react-bootstrap/react-bootstrap.js",
    "bower_components/react-router-bootstrap/lib/ReactRouterBootstrap.js",
    "app/scripts/**/*.jsx"
  ])
    .pipe(order([
    "bower_components/jquery/dist/jquery.js",
    "bower_components/react/JSXTransformer.js",
    "bower_components/react/react-with-addons.js",
    "bower_components/react-router/build/umd/ReactRouter.js",
    "bower_components/react-bootstrap/react-bootstrap.js",
    "bower_components/react-bootstrap/react-bootstrap.js",
    "bower_components/react-router-bootstrap/lib/ReactRouterBootstrap.js",
    "app/scripts/app.jsx",
    "app/scripts/components/*",
    "app/scripts/actions/*",
    "app/scripts/constants/*",
    "app/scripts/stores/*",
    "app/scripts/dispatcher/*",
    "app/scripts/routes.jsx"
      ], { base: './' }))
    .pipe(gulpif(args.sourcemaps, sourcemaps.init()))
    .pipe(concat("application.js"))
    .pipe(react())
    .pipe(gulpif(args.sourcemaps, sourcemaps.write()))
    .pipe(gulpif(args.production, uglify()))
    .pipe(gulp.dest("public/static/"))
    .pipe(reload());
});
 
/*
 * Create a customized modernizr build.
 */
gulp.task("app:modernizr", function() {
  return gulp.src([
    //"public/stylesheets/style.css",
    "public/javascripts/application.js"
  ]).pipe(
      modernizr({
        options: [
          "addTest",                   /* Add custom tests */
          "fnBind",                    /* Use function.bind */
          "html5printshiv",            /* HTML5 support for IE */
          "setClasses",                /* Add CSS classes to root tag */
          "testProp"                   /* Test for properties */
        ]
      }))
    .pipe(addsrc.append("bower_components/respond/dest/respond.src.js"))
    .pipe(concat("modernizr.js"))
    .pipe(gulpif(args.production, uglify()))
    .pipe(gulp.dest("public/static"));
});
 
/*
 * Minify views.
 */
gulp.task("app:views", args.production ? [
  "app:stylesheets",
  "app:javascripts",
  "app:modernizr",
  "app:revisions:clean",
  "app:revisions"
] : [], function() {
  return gulp.src(["views/**/*.tmpl"])
    .pipe(cache("views"))
    .pipe(
      minhtml({
        collapseBooleanAttributes: true,
        collapseWhitespace: true,
        removeComments: true,
        removeScriptTypeAttributes: true,
        removeStyleLinkTypeAttributes: true,
        minifyCSS: true,
        minifyJS: true
      }))
    .pipe(gulp.dest(".views"));
});
 
/*
 * Clean outdated revisions.
 */
gulp.task("app:revisions:clean", function() {
  return gulp.src(["public/**/*.{css,js}"])
    .pipe(ignore.include(/-[a-f0-9]{8}\.(css|js)$/))
    .pipe(vinyl(clean));
});
 
/*
 * Revision app after build.
 */
gulp.task("app:revisions", [
  "app:stylesheets",
  "app:javascripts",
  "app:modernizr",
  "app:revisions:clean"
], function() {
  return gulp.src(["public/**/*.{css,js}"])
    .pipe(ignore.exclude(/-[a-f0-9]{8}\.(css|js)$/))
    .pipe(rev())
    .pipe(gulp.dest("public"))
    .pipe(rev.manifest())
    .pipe(gulp.dest("public"));
})
 
/*
 * Build app.
 */
gulp.task("app:build", [
  //"app:stylesheets",
  "app:javascripts",
  "app:modernizr"
]);
 
/*
 * Watch app for changes and rebuild on the fly.
 */
gulp.task("app:watch", function() {
  
  /* Rebuild stylesheets on-the-fly */
  //gulp.watch([
  //  "app/styles/**/*.scss"
  //], ["app:stylesheets"]);
 
  /* Rebuild javascripts on-the-fly */
  gulp.watch([
    "app/scripts/**/*.js",
    "app/scripts/**/*.jsx",
    "bower.json"
  ], ["app:javascripts"]);
});
 
/* ----------------------------------------------------------------------------
 * Application server
 * ------------------------------------------------------------------------- */
 
/*
 * Build application server.
 */
gulp.task("server:build", function() {
  if (server)
    server.kill();
  var build = child.spawnSync("go", ["install"]);
  if (build.stderr.length) {
    var lines = build.stderr.toString()
      .split("\n").filter(function(line) {
        return line.length
      });
    for (var l in lines)
      util.log(util.colors.red(
        "Error (go install): " + lines[l]
      ));
    notifier.notify({
      title: "Error (go install)",
      message: lines
    });
  }
  return build;
});
 
/*
 * Restart application server.
 */
gulp.task("server:spawn", function() {
  if (server)
    server.kill();
 
  /* Spawn application server */
  server = child.spawn("tyovuoro.exe");
 
  /* Trigger reload upon server start */
  server.stdout.once("data", function() {
    reload.reload("/");
  });
 
  /* Pretty print server log output */
  server.stdout.on("data", function(data) {
    var lines = data.toString().split("\n")
    for (var l in lines)
      if (lines[l].length)
        util.log(lines[l]);
  });
 
  /* Print errors to stdout */
  server.stderr.on("data", function(data) {
    process.stdout.write(data.toString());
  });
});
 
/*
 * Watch source for changes and restart application server.
 */
gulp.task("server:watch", function() {
 
  /* Restart application server */
  gulp.watch([
    ".views/**/*.tmpl",
    "locales/*.json"
  ], ["server:spawn"]);
 
  /* Rebuild and restart application server */
  gulp.watch([
    "*.go",
    "*/**/*.go"
  ], sync([
    "server:build",
    "server:spawn"
  ], "server"));
});
 
/* ----------------------------------------------------------------------------
 * Interface
 * ------------------------------------------------------------------------- */
 
/*
 * Build assets and application server.
 */
gulp.task("build", [
  "app:build",
  "server:build"
]);
 
/*
 * Start asset and server watchdogs and initialize livereload.
 */
gulp.task("watch", [
/*  "app:build",*/
  "server:build"
], function() {
  reload.listen();
  return gulp.start([
    "app:build",
    "app:watch",
    "server:watch",
    "server:spawn"
  ]);
});
 
/*
 * Build app by default.
 */
gulp.task("default", ["build"]);