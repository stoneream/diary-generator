ThisBuild / scalaVersion := "2.13.8"
ThisBuild / version := "0.2.0-SNAPSHOT"
ThisBuild / organization := "com.github.stoneream"

lazy val root = (project in file(".")).settings(
  name := "diary-generator",
  assembly / mainClass := Some("app.Main"),
  assembly / assemblyJarName := "diary-generator.jar",
  semanticdbEnabled := true,
  semanticdbVersion := scalafixSemanticdb.revision,
  libraryDependencies ++= Seq(
    "ch.qos.logback" % "logback-classic" % "1.2.10",
    "com.typesafe.scala-logging" %% "scala-logging" % "3.9.5",
    "com.github.nscala-time" %% "nscala-time" % "2.32.0",
    "com.github.scopt" %% "scopt" % "4.1.0",
    "org.scalatest" %% "scalatest" % "3.2.8"
  )
)

scalacOptions ++= List(
  "-Ywarn-unused",
  "-Yrangepos"
)

scalafmtOnCompile := true
