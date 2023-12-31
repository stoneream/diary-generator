import sbt._

object Dependencies {

  object Versions {
    val logback = "1.4.7"
    val scalaLogging = "3.9.5"
    val nscalatime = "2.32.0"
    val scopt = "4.1.0"
    val scalatest = "3.2.15"
  }

  lazy val bundle: Seq[ModuleID] = logging ++ nscalatime ++ scopt ++ scalatest

  lazy val logging: Seq[ModuleID] = Seq(
    "ch.qos.logback" % "logback-classic" % Versions.logback,
    "com.typesafe.scala-logging" %% "scala-logging" % Versions.scalaLogging
  )

  lazy val nscalatime: Seq[ModuleID] = Seq(
    "com.github.nscala-time" %% "nscala-time" % Versions.nscalatime
  )

  lazy val scopt: Seq[ModuleID] = Seq(
    "com.github.scopt" %% "scopt" % Versions.scopt
  )

  lazy val scalatest: Seq[ModuleID] = Seq(
    "org.scalatest" %% "scalatest" % Versions.scalatest % Test
  )

}
